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
	query := fmt.Sprintf(`INSERT INTO %s (first_name, last_name, username, password_hash, rating, num_people_rated, tg_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`, usersTable)
	row := r.db.QueryRow(query, user.FirstName, user.LastName, user.Username, user.Password, user.Rating, user.NumPeopleRated, user.TgId)
	if err := row.Scan(&id); err != nil {
		return -1, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUserId(username, passwordHash string) (int, error) {
	var id int
	query := fmt.Sprintf(`SELECT id FROM %s WHERE username=$1 and password_hash=$2`, usersTable)
	err := r.db.Get(&id, query, username, passwordHash)
	return id, err
}

func (r *AuthPostgres) GetTgUser(username string, tgId int) (int, error) {
	var userId int
	query := fmt.Sprintf(`SELECT id  FROM %s WHERE username=$1 and tg_id=$2`, usersTable)
	err := r.db.Get(&userId, query, username, tgId)

	return userId, err
}
