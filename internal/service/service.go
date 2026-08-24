package service

import (
	"context"
	"errors"
	"hamrahTask1/internal/core/ports"
	"hamrahTask1/pkg/logger"

	"github.com/google/uuid"
)

type chatService struct {
	userRepo ports.UserRepository
	chatRepo ports.ChatRepository
	cache    ports.SessionCache
	log      *logger.Logger
}

func NewChatService(ur ports.UserRepository, cr ports.ChatRepository, sc ports.SessionCache, l *logger.Logger) ports.Service {
	return &chatService{userRepo: ur, chatRepo: cr, cache: sc, log: l}
}

func (s *chatService) Register(ctx context.Context, username, password string, age int64, gender bool) error {
	s.log.Info("Registering user: " + username)
	return s.userRepo.CreateUser(ctx, username, password, age, gender)
}

func (s *chatService) Login(ctx context.Context, username, password string) (string, error) {
	id, hash, err := s.userRepo.GetUserByUsername(ctx, username)
	if err != nil || hash != password {
		s.log.Error("Failed login attempt for", err)
		return "", errors.New("invalid credentials")
	}

	s.userRepo.UpdateOnlineStatus(ctx, username, true, true)
	token := uuid.New().String()
	err = s.cache.CreateSession(ctx, token, id)
	return token, err
}

func (s *chatService) Logout(ctx context.Context, token string) error {
	id, err := s.cache.GetUserID(ctx, token)
	if err == nil {
		s.userRepo.UpdateOnlineStatus(ctx, id, false, false)
		s.cache.DeleteSession(ctx, token)
	}
	return err
}

func (s *chatService) GetUsers(ctx context.Context) ([]string, []bool, error) {
	return s.userRepo.GetAllUsers(ctx)
}

func (s *chatService) SendMessage(ctx context.Context, token, receiverUsername, message string) error {
	senderID, err := s.cache.GetUserID(ctx, token)
	if err != nil {
		return errors.New("unauthorized")
	}
	receiverID, err := s.chatRepo.GetIDByUsername(ctx, receiverUsername)
	if err != nil {
		return errors.New("receiver not found")
	}
	return s.chatRepo.SaveMessage(ctx, senderID, receiverID, message)
}

func (s *chatService) GetMessages(ctx context.Context, token, contactUsername string, count int64) ([]string, error) {
	userID, err := s.cache.GetUserID(ctx, token)
	if err != nil {
		return nil, errors.New("unauthorized")
	}
	contactID, err := s.chatRepo.GetIDByUsername(ctx, contactUsername)
	if err != nil {
		return nil, errors.New("contact not found")
	}
	return s.chatRepo.GetMessages(ctx, userID, contactID, count)
}
