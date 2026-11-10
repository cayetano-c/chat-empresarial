package services

import (
	"github.com/cayetano-c/chat-empresarial/app-go/internal/auth"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/models"
	"github.com/cayetano-c/chat-empresarial/app-go/internal/storage"
)

const PageSize = 50

type ChatService struct {
	store storage.MessageStore
	users auth.Authenticator
}

func NewChatService(store storage.MessageStore, users auth.Authenticator) *ChatService {
	return &ChatService{store: store, users: users}
}

func (s *ChatService) SendMessage(from, to, text string) (*models.Message, error) {
	msg, err := models.NewMessage(from, to, text)
	if err != nil {
		return nil, err 
	}
	if !s.users.UserExists(msg.To) {
		return nil, ErrUserNotFound
	}
	if err := s.store.Save(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *ChatService) NewMessages(user string, afterID int) []models.Message {
	return s.store.Pending(user, afterID)
}

func (s *ChatService) History(user, with string, page int) ([]models.Message, error) {
	if page < 1 {
		return nil, ErrInvalidPage
	}
	if !s.users.UserExists(with) {
		return nil, ErrUserNotFound
	}

	all := s.store.Conversation(user, with)
	start := (page - 1) * PageSize
	if start >= len(all) {
		return []models.Message{}, nil // no hay más mensajes
	}
	end := start + PageSize
	if end > len(all) {
		end = len(all)
	}
	return all[start:end], nil
}
