package repository

import (
	"fmt"

	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetUserInfo(id int64) (core.User, error) {
	var destUser core.User
	query := fmt.Sprintf(`SELECT id, first_name, last_name, username, rating, num_people_rated, tg_id  FROM %s WHERE id=$1`, usersTable)
	err := r.db.Get(&destUser, query, id)
	return destUser, err
}

func (r *UserPostgres) JoinTrip(userId, tripId int64) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`SELECT places_taken, places_max FROM %s WHERE id=$1`, tripsTable)
	var trip core.Trip
	err = tx.QueryRow(query, tripId).Scan(&trip.PlacesTaken, &trip.PlacesMax)
	if err != nil {
		tx.Rollback()
		return err
	}

	if trip.PlacesTaken+1 > trip.PlacesMax {
		tx.Rollback()
		return fmt.Errorf("no more places left")
	}

	query = fmt.Sprintf(`INSERT INTO %s (user_id, trip_id) VALUES ($1, $2)`, usersTripsTable)
	_, err = tx.Exec(query, userId, tripId)
	if err != nil {
		tx.Rollback()
		return err
	}

	query = fmt.Sprintf(`UPDATE %s SET places_taken=places_taken+1 WHERE id=$1`, tripsTable)
	_, err = tx.Exec(query, tripId)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}
