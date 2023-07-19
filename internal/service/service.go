package service

import (
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/internal/repository"
)

type Authorization interface {
	CreateUser(user core.User) (int64, error)
	GetUserId(user core.User) (int64, error)

	// GenerateToken(userId, tgId int) (string, error)
	// ParseToken(accessToken string) (int, error)
	GenerateToken(id core.UserCtx) (string, error)
	ParseToken(accessToken string) (core.UserCtx, error)

	VerifyTgAuthData(authData map[string]interface{}, keyword string) (bool, error) // WARNING! only error or bool should be returned
}

type User interface {
	GetUserInfo(id int64) (core.User, error)
	JoinTrip(userId, tripId int64) error
	// RateUser(core.User) (int, error)
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

type Service struct {
	Authorization
	User
	Trip
	Message
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		User:          NewUserService(repo.User),
		Trip:          NewTripService(repo.Trip),
		Message:       NewMessageService(repo.Message),
	}
}
