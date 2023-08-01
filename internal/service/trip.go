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

func (s *TripService) Create(trip core.Trip) (int64, error) {
	return s.repo.Create(trip)
}

func (s *TripService) GetById(tripId int64) (core.Trip, error) {
	return s.repo.GetById(tripId)
}

func (s *TripService) Delete(userId, tripId int64) (int64, error) {
	return s.repo.Delete(userId, tripId)
}

func (s *TripService) GetAdjTrips(input core.InputAdjTrips) ([]core.Trip, error) {
	return s.repo.GetAdjTrips(input)
}

func (s *TripService) GetJoinedTrips(userId int64) ([]core.Trip, error) {
	return s.repo.GetJoinedTrips(userId)
}

func (s *TripService) GetJoinedUsers(userId, tripId int64) ([]core.UserCtx, error) {
	return s.repo.GetJoinedUsers(userId, tripId)
}
