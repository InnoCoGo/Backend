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
	query := fmt.Sprintf(`INSERT INTO %s (user_id, trip_id) VALUES ($1, $2)`, usersTripsTable)
	_, err := r.db.Exec(query, userId, tripId)
	return err
}
