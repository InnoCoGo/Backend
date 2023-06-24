package repository

import (
	"fmt"

	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
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
	row := r.db.QueryRow(query, user.FirstName, user.LastName, user.Username, user.PasswordOrHash, user.Rating, user.NumPeopleRated, user.TgId)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	// log.Printf("\nCreateUser REPOS: %v\n", id)
	return id, nil
}

type getUserIdRes struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
}

func (r *AuthPostgres) GetUserId(user core.User) (int, error) {
	var id int

	var query string
	var err error
	var userRes getUserIdRes

	if user.PasswordOrHash != "" {
		query = fmt.Sprintf(`SELECT id FROM %s WHERE username=$1 and password_hash=$2`, usersTable)
		err = r.db.Get(&id, query, user.Username, user.PasswordOrHash)
	} else {
		query = fmt.Sprintf(`SELECT id, username FROM %s WHERE tg_id=$1`, usersTable)
		err = r.db.Get(&userRes, query, user.TgId)
		if err == nil { // WARNING! UPDATE can be done through another interface
			id = userRes.Id
			if userRes.Username != user.Username {
				query = fmt.Sprintf(`UPDATE %s SET username=$1 WHERE tg_id=$2`, usersTable)
				_, err = r.db.Exec(query, user.Username, user.TgId)
			}
		}
	}
	return id, err
}
