package service

import (
	"github.com/itoqsky/InnoCoTravel_backend/internal/core"
	"github.com/itoqsky/InnoCoTravel_backend/internal/repository"
)

type Authorization interface {
	CreateUser(user core.User) (int, error)
	GenerateToken(username, password string) (string, error)
}

// type User interface {
// 	RateUser(core.User) (int, error)
// }

// type Trip interface {
// 	Create(trip core.Trip) (int, error)
// }

type Service struct {
	Authorization
	// User
	// Trip
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.Authorization),
	}
}
