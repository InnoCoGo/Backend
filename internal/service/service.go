package service

import (
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/internal/repository"
)

type Authorization interface {
	CreateUser(user core.User) (int, error)
	GetUserId(user core.User) (int, error)

	// GenerateToken(userId, tgId int) (string, error)
	// ParseToken(accessToken string) (int, error)
	GenerateToken(id core.UserCtx) (string, error)
	ParseToken(accessToken string) (core.UserCtx, error)

	VerifyTgAuthData(authData map[string]interface{}, keyword string) (bool, error) // WARNING! only error or bool should be returned
}

type User interface {
	GetUserInfo(id int) (core.User, error)
	// RateUser(core.User) (int, error)
}

type Trip interface {
	Create(trip core.Trip) (int, error)
	GetById(userId, tripId int) (core.Trip, error)
	Delete(userId, tripId int) (int, error)
	// Update(trip core.Trip) error
	GetAdjTrips(input core.InputAdjTrips) ([]core.Trip, error)
	GetJoinedTrips(userId int) ([]core.Trip, error)
}

type Service struct {
	Authorization
	User
	Trip
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
		User:          NewUserService(repo.User),
		Trip:          NewTripService(repo.Trip),
	}
}
