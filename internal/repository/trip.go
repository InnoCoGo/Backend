package repository

import (
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
