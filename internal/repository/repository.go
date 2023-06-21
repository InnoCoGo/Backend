package repository

import (
	"github.com/itoqsky/InnoCoTravel_backend/internal/core"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user core.User) (int, error)
	GetUserId(user core.User) (int, error)
	// LoginTg(core.User) (string, error)
}
type User interface {
	GetUserInfo(id int) (core.User, error)
}

type Trip interface {
	Create(trip core.Trip) (int, error)
	GetById(userId, tripId int) (core.Trip, error)
	Delete(userId, tripId int) (int, error)
	// Update(trip core.Trip) error
	GetAdjTrips(input core.InputAdjTrips) ([]core.Trip, error)
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
