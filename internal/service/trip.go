package service

import (
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/internal/repository"
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

func (s *TripService) GetById(userId, tripId int) (core.Trip, error) {
	return s.repo.GetById(userId, tripId)
}

func (s *TripService) Delete(userId, tripId int) (int, error) {
	return s.repo.Delete(userId, tripId)
}

// func (s *TripService) Update(trip core.Trip) error {
// 	return s.repo.Update(trip)
// }

func (s *TripService) GetAdjTrips(input core.InputAdjTrips) ([]core.Trip, error) {
	return s.repo.GetAdjTrips(input)
}

func (s *TripService) GetJoinedTrips(userId int) ([]core.Trip, error) {
	return s.repo.GetJoinedTrips(userId)
}
