package service

import (
	"github.com/itoqsky/InnoCoTravel_backend/internal/core"
	"github.com/itoqsky/InnoCoTravel_backend/internal/repository"
)

type Authorization interface {
	CreateUser(user core.User) (int, error)

	GenerateToken(username, password string) (string, error)
	ParseToken(accessToken string) (int, error)

	GetTgUser(user core.User) (int, error)
}

type User interface {
	GetUserInfo(id int) (core.User, error)
	// RateUser(core.User) (int, error)
}

type Trip interface {
	Create(trip core.Trip) (int, error)
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
