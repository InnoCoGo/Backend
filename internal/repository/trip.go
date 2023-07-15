package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/jmoiron/sqlx"
)

type TripPostgres struct {
	db *sqlx.DB
}

func NewTripPostgres(db *sqlx.DB) *TripPostgres {
	return &TripPostgres{db: db}
}

func (r *TripPostgres) Create(trip core.Trip) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int64
	createTripQuery := fmt.Sprintf(`INSERT INTO %s (admin_id, admin_username, admin_tg_id, is_driver, places_max, places_taken, chosen_timestamp, from_point, to_point, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`, tripsTable)
	row := tx.QueryRow(createTripQuery, trip.AdminId, trip.AdminUsername, trip.AdminTgId, trip.IsDriver, trip.PlacesMax, trip.PlacesTaken, trip.ChosenTimestamp, trip.FromPoint, trip.ToPoint, trip.Description)
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

func (r *TripPostgres) GetById(tripId int64) (core.Trip, error) {
	var trip core.Trip
	query := fmt.Sprintf(`SELECT * FROM %s t WHERE t.id=$1`, tripsTable)
	err := r.db.Get(&trip, query, tripId)
	return trip, err
}

func (r *TripPostgres) Delete(userId, tripId int64) (int64, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var adminId, adminTgId int64
	var adminUsername string
	var placesTaken int
	getTripQuery := fmt.Sprintf(`SELECT t.admin_id, t.admin_username, t.admin_tg_id, t.places_taken FROM %s t INNER JOIN %s ut ON ut.trip_id=t.id WHERE ut.trip_id=$1 AND ut.user_id=$2`, tripsTable, usersTripsTable)
	row := tx.QueryRow(getTripQuery, tripId, userId)
	if err := row.Scan(&adminId, &adminUsername, &adminTgId, &placesTaken); err != nil {
		tx.Rollback()
		return 0, err
	}

	deleteUserQuery := fmt.Sprintf(`DELETE FROM %s ut WHERE ut.user_id=$1 AND ut.trip_id=$2`, usersTripsTable)
	_, err = tx.Exec(deleteUserQuery, userId, tripId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if adminId == userId {
		newAdminId, newAdminUsername, newAdminTgId, err := r.chooseNewAdmin(tx, tripId)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
		if newAdminId == 0 {
			return 0, tx.Commit()
		}
		adminId = newAdminId
		adminUsername = newAdminUsername
		adminTgId = newAdminTgId
	}
	placesTaken--

	setValues := `admin_id=$1, admin_username=$2, admin_tg_id=$3, places_taken=$4`
	deleteTripQuery := fmt.Sprintf(`UPDATE %s SET %s WHERE id=$5`, tripsTable, setValues)
	_, err = tx.Exec(deleteTripQuery, adminId, adminUsername, adminTgId, placesTaken, tripId)
	if err != nil {
		log.Printf("\nAAAA\n")
		tx.Rollback()
		return 0, err
	}

	// log.Printf("\n %v | %v | %v \n", adminId, adminUsername, adminTgId)

	return adminId, tx.Commit()
}

func (r *TripPostgres) chooseNewAdmin(tx *sql.Tx, tripId int64) (newAdminId int64, newAdminUsername string, newAdminTgId int64, err error) {
	nextAdminQuery := fmt.Sprintf(`SELECT u.id, u.username, u.tg_id 
								FROM 
									%s u 
								INNER JOIN 
									%s ut
								ON ut.user_id=u.id
								WHERE
									ut.trip_id=$1`, usersTable, usersTripsTable)
	row := tx.QueryRow(nextAdminQuery, tripId)

	if err = row.Scan(&newAdminId, &newAdminUsername, &newAdminTgId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tripQuery := fmt.Sprintf(`DELETE FROM %s t WHERE t.id=$1`, tripsTable)
			_, err = tx.Exec(tripQuery, tripId)
			if err != nil {
				return 0, "", 0, err
			}
		} else {
			return 0, "", 0, err
		}
	}
	return newAdminId, newAdminUsername, newAdminTgId, nil
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

func (r *TripPostgres) GetJoinedTrips(userId int64) ([]core.Trip, error) {
	var dest []core.Trip
	query := fmt.Sprintf(`SELECT t.*
						FROM 
							%s as t
						INNER JOIN %s as ut 
							ON  ut.trip_id = t.id
						WHERE 
							ut.user_id = $1`, tripsTable, usersTripsTable)
	err := r.db.Select(&dest, query, userId)
	return dest, err
}

func (r *TripPostgres) GetJoinedUsers(userId, tripId int64) ([]core.UserCtx, error) {
	var dest []core.UserCtx
	query := fmt.Sprintf(`SELECT u.id, u.username
						FROM 
							%s u
						INNER JOIN %s ut
							ON  ut.user_id = u.id
						WHERE ut.user_id=$1
							AND ut.trip_id=$2
	`, usersTable, usersTripsTable)
	err := r.db.Select(&dest, query, userId, tripId)

	return dest, err
}
