package handler

import (
	"context"
	"hamrahTask1/internal/core/ports"
	"hamrahTask1/pkg/logger"
)

type GRPCHandler struct {
	UnimplementedChatbotServiceServer
	svc ports.Service
	log *logger.Logger
}

func NewGRPCHandler(svc ports.Service, log *logger.Logger) *GRPCHandler {
	return &GRPCHandler{svc: svc, log: log}
}

func (h *GRPCHandler) Register(ctx context.Context, req *RegisterRequest) (*ReportResponse, error) {
	err := h.svc.Register(ctx, req.GetUsername(), req.GetPassword(), req.GetAge(), req.GetGender())
	if err != nil {
		return &ReportResponse{Message: "Registration failed", Success: false}, nil
	}
	return &ReportResponse{Message: "registered successfully!", Success: true}, nil
}

func (h *GRPCHandler) Login(ctx context.Context, req *LoginRequest) (*ReportResponse, error) {
	token, err := h.svc.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return &ReportResponse{Message: "Invalid credentials", Success: false}, nil
	}
	return &ReportResponse{Message: "Login successful", Success: true, Id: token}, nil
}

func (h *GRPCHandler) Logout(ctx context.Context, req *LogoutRequest) (*ReportResponse, error) {
	h.svc.Logout(ctx, req.GetId())
	return &ReportResponse{Message: "Logged out successful", Success: true}, nil
}

func (h *GRPCHandler) ShowUserList(ctx context.Context, req *Empty) (*UserListResponse, error) {
	users, status, err := h.svc.GetUsers(ctx)
	if err != nil {
		return nil, err
	}
	return &UserListResponse{UserList: users, Status: status}, nil
}

func (h *GRPCHandler) SendMessage(ctx context.Context, req *SendMessageRequest) (*ReportResponse, error) {
	err := h.svc.SendMessage(ctx, req.GetSenderId(), req.GetReceiverId(), req.GetMessage())
	if err != nil {
		return &ReportResponse{Message: err.Error(), Success: false}, nil
	}
	return &ReportResponse{Message: "Message sent!", Success: true}, nil
}

func (h *GRPCHandler) ShowMessage(ctx context.Context, req *ShowMessagesRequest) (*ShowMessagesResponse, error) {
	msgs, err := h.svc.GetMessages(ctx, req.GetId(), req.GetContactId(), req.GetCount())
	if err != nil {
		return &ShowMessagesResponse{Message: []string{err.Error()}, Success: false}, nil
	}
	return &ShowMessagesResponse{Message: msgs, Success: true}, nil
}
