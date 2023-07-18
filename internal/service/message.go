package service

import (
	"github.com/itoqsky/InnoCoTravel-backend/internal/core"
	"github.com/itoqsky/InnoCoTravel-backend/internal/repository"
)

type MessageService struct {
	repo repository.Message
}

func NewMessageService(repo repository.Message) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) Save(message core.Message) (int64, error) {
	return s.repo.Save(message)
}

func (s *MessageService) FetchRoomMessages(roomId int64) ([]core.Message, error) {
	return s.repo.FetchRoomMessages(roomId)
}
