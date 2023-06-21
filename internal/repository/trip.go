package repository

import (
	"database/sql"
	"fmt"

	"github.com/itoqsky/InnoCoTravel_backend/internal/core"
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
		return -1, err
	}

	var id int
	createTripQuery := fmt.Sprintf(`INSERT INTO %s (admin_id, is_passanger, places_max, places_taken, chosen_date_time, from_point, to_point, description)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`, tripsTable)
	row := tx.QueryRow(createTripQuery, trip.AdminId, trip.IsPassanger, trip.PlacesMax, trip.PlacesTaken, trip.ChosenDateTime, trip.FromPoint, trip.ToPoint, trip.Description)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return -1, err
	}

	createUsersTripQuery := fmt.Sprintf(`INSERT INTO %s (user_id, trip_id) VALUES ($1, $2)`, usersTripsTable)
	_, err = tx.Exec(createUsersTripQuery, trip.AdminId, id)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return id, tx.Commit()
}

func (r *TripPostgres) GetById(userId, tripId int) (core.Trip, error) {
	query := fmt.Sprintf(`SELECT
							t.id,
							t.admin_id,
							t.is_passanger,
							t.places_max,
							t.places_taken,
							t.chosen_date_time,
							t.from_point,
							t.to_point,
							t.description
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

func (r *TripPostgres) Delete(userId, tripId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	usersTripsQuery := fmt.Sprintf(`DELETE FROM %s ut WHERE ut.user_id=$1 AND ut.trip_id=$2`, usersTripsTable)
	_, err = tx.Exec(usersTripsQuery, userId, tripId)
	if err != nil {
		tx.Rollback()
		return err
	}

	var nextAdmin int
	nextAdminQuery := fmt.Sprintf(`SELECT ut.user_id FROM %s ut WHERE ut.trip_id=$1`, usersTripsQuery)
	row := tx.QueryRow(nextAdminQuery, tripId)
	if err := row.Scan(&nextAdmin); err != nil {
		if err == sql.ErrNoRows {
			tripQuery := fmt.Sprintf(`DELETE FROM %s t WHERE t.id=$1 AND t.places_taken=0`, tripsTable)
			_, err = tx.Exec(tripQuery, tripId)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func (r *TripPostgres) Update(trip core.Trip) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

}
