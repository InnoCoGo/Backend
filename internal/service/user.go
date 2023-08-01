package service

import (
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/internal/repository"
)

type UserService struct {
	repo repository.User
}

func NewUserService(repo repository.User) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserInfo(id int64) (core.User, error) {
	return s.repo.GetUserInfo(id)
}

func (s *UserService) JoinTrip(userId, tripId int64) error {
	return s.repo.JoinTrip(userId, tripId)
}
