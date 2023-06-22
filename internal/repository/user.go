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

func (r *UserPostgres) GetUserInfo(id int) (core.User, error) {
	var destUser core.User
	query := fmt.Sprintf(`SELECT id, first_name, last_name, username, rating, num_people_rated, tg_id  FROM %s WHERE id=$1`, usersTable)
	err := r.db.Get(&destUser, query, id)
	return destUser, err
}
