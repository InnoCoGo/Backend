package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/jmoiron/sqlx"
)

type TripPostgres struct {
	db *sqlx.DB
}

func NewTripPostgres(db *sqlx.DB) *TripPostgres {
	return &TripPostgres{db: db}
}

func (r *TripPostgres) Create(trip core.Trip) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createTripQuery := fmt.Sprintf(`INSERT INTO %s (admin_id, admin_username, is_driver, places_max, places_taken, chosen_timestamp, from_point, to_point, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`, tripsTable)
	row := tx.QueryRow(createTripQuery, trip.AdminId, trip.AdminUsername, trip.IsDriver, trip.PlacesMax, trip.PlacesTaken, trip.ChosenTimestamp, trip.FromPoint, trip.ToPoint, trip.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	// log.Printf("\n%v\n", trip.AdminId)

	createUsersTripQuery := fmt.Sprintf(`INSERT INTO %s (user_id, trip_id) VALUES ($1, $2)`, usersTripsTable)
	_, err = tx.Exec(createUsersTripQuery, trip.AdminId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return id, tx.Commit()
}

func (r *TripPostgres) GetById(userId, tripId int) (core.Trip, error) {
	query := fmt.Sprintf(`SELECT *
						FROM 
							%s t
						INNER JOIN %s ut
							ON  ut.trip_id = t.id
						WHERE ut.user_id=$1
							AND ut.trip_id=$2
	`, tripsTable, usersTripsTable)
	var trip core.Trip
	err := r.db.Get(&trip, query, userId, tripId)

	return trip, err
}

func (r *TripPostgres) Delete(userId, tripId int) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	usersTripsQuery := fmt.Sprintf(`DELETE FROM %s ut WHERE ut.user_id=$1 AND ut.trip_id=$2`, usersTripsTable)
	_, err = tx.Exec(usersTripsQuery, userId, tripId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var newAdminId int
	nextAdminQuery := fmt.Sprintf(`SELECT ut.user_id FROM %s ut WHERE ut.trip_id=$1`, usersTripsTable)
	row := tx.QueryRow(nextAdminQuery, tripId)

	if err := row.Scan(&newAdminId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tripQuery := fmt.Sprintf(`DELETE FROM %s t WHERE t.id=$1 AND t.places_taken=0`, tripsTable)
			_, err = tx.Exec(tripQuery, tripId)
			if err != nil {
				tx.Rollback()
				return 0, err
			}
		} else {
			tx.Rollback()
			return 0, err
		}
	} else {
		setValues := `admin_id=$1`
		nextAdminQuery = fmt.Sprintf(`UPDATE %s SET %s WHERE t.id=$2`, tripsTable, setValues)
		_, err = tx.Exec(nextAdminQuery, newAdminId)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	return newAdminId, tx.Commit()
}

// func (r *TripPostgres) Update(trip core.Trip) error {
// 	tx, err := r.db.Begin()
// 	if err != nil {
// 		return err
// 	}

// }

func (r *TripPostgres) GetAdjTrips(input core.InputAdjTrips) ([]core.Trip, error) {
	var trips []core.Trip
	companionValues := ""
	args := make([]interface{}, 0)

	if input.CompanionType == "passenger" {
		companionValues = "AND is_driver=FALSE"
	} else if input.CompanionType == "driver" {
		companionValues = "AND is_driver=TRUE"
	}

	args = append(args, input.LeftTimestamp, input.RightTimestamp, input.FromPoint, input.ToPoint)

	query := fmt.Sprintf(`SELECT * FROM %s 
		WHERE $1 <= chosen_timestamp AND chosen_timestamp <= $2 AND from_point=$3 AND to_point=$4 %s`, tripsTable, companionValues)

	err := r.db.Select(&trips, query, args...)

	return trips, err
}

func (r *TripPostgres) GetJoinedTrips(userId int) ([]core.Trip, error) {
	var dest []core.Trip
	query := fmt.Sprintf(`SELECT 
							t.id,
							t.admin_id,
							t.admin_username,
							t.is_driver,
							t.places_max,
							t.places_taken,
							t.chosen_timestamp,
							t.from_point,
							t.to_point,
							t.description
						FROM %s as t, %s as u WHERE u.user_id = $1`, tripsTable, usersTripsTable)
	err := r.db.Select(&dest, query, userId)
	return dest, err
}
