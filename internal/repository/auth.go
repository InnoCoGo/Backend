package repository

import (
	"fmt"

	"github.com/itoqsky/InnoCoTravel_backend/internal/core"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user core.User) (int, error) {
	var id int
	query := fmt.Sprintf(`INSERT INTO %s (f_name, l_name, username, password_hash, rating, num_people_rated) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`, usersTable)
	row := r.db.QueryRow(query, user.FirstName, user.LastName, user.Username, user.Password, user.Rating, user.NumPeopleRated)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, passwordHash string) (core.User, error) {
	var destUser core.User
	query := fmt.Sprintf(`SELECT id, f_name, l_name, username, rating, num_people_rated  FROM %s WHERE username=$1 and password_hash=$2`, usersTable)
	err := r.db.Get(&destUser, query, username, passwordHash)
	return destUser, err
}
