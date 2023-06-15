package repository

import (
	"github.com/itoqsky/InnoCoTravel_backend/internal/core"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user core.User) (int, error)
	GetUserId(username, passwordHash string) (int, error)
	GetTgUser(username string, tgId int) (int, error)
	// LoginTg(core.User) (string, error)
}
type User interface {
	GetUserInfo(id int) (core.User, error)
}

type Trip interface {
	Create(trip core.Trip) (int, error)
}

type Repository struct {
	Authorization
	User
	Trip
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		User:          NewUserPostgres(db),
		Trip:          NewTripPostgres(db),
	}
}
