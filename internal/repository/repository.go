package repository

import (
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/jmoiron/sqlx"
)

type Authorization interface {
	CreateUser(user core.User) (int64, error)
	GetUserId(user core.User) (int64, error)
	// LoginTg(core.User) (string, error)
}
type User interface {
	GetUserInfo(id int64) (core.User, error)
	JoinTrip(userId, tripId int64) error
}

type Trip interface {
	Create(trip core.Trip) (int64, error)
	GetById(userId, tripId int64) (core.Trip, error)
	Delete(userId, tripId int64) (int64, error)
	GetAdjTrips(input core.InputAdjTrips) ([]core.Trip, error)
	GetJoinedTrips(userId int64) ([]core.Trip, error)
	GetJoinedUsers(userId, tripId int64) ([]core.UserCtx, error)
}

type Message interface {
	Save(message core.Message) (int64, error)
	FetchRoomMessages(roomId int64) ([]core.Message, error)
}

type Repository struct {
	Authorization
	User
	Trip
	Message
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		User:          NewUserPostgres(db),
		Trip:          NewTripPostgres(db),
		Message:       NewMessagePostgres(db),
	}
}
