package service

import (
	"github.com/itoqsky/InnoCoTravel_backend/internal/core"
	"github.com/itoqsky/InnoCoTravel_backend/internal/repository"
)

type TripService struct {
	repo repository.Trip
}

func NewTripService(repo repository.Trip) *TripService {
	return &TripService{repo: repo}
}

func (s *TripService) Create(trip core.Trip) (int, error) {
	return s.repo.Create(trip)
}
