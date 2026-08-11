package main

import (
	"context"
	"database/sql"
	"fmt"
	"hamrahTask1/proto"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type server struct {
	proto.UnimplementedChatbotServiceServer
	db  *sql.DB
	rdb *redis.Client
}

func (s *server) Register(ctx context.Context, req *proto.RegisterRequest) (*proto.ReportResponse, error) {

	username := req.GetUsername()
	password := req.GetPassword()
	gender := req.GetGender()
	age := req.GetAge()

	_, err := s.db.ExecContext(ctx,
		"INSERT INTO users (username, password_hash, age, gender, is_online) VALUES ($1, $2, $3, $4, $5)",
		username, password, age, gender, false)

	if err != nil {
		log.Printf("Warning: Failed to save user to database: %v", err)
		return &proto.ReportResponse{
			Message: "Registration failed (user may already exist)",
			Success: false,
		}, nil
	}

	log.Printf("Registering user with username %s", username)
	return &proto.ReportResponse{
		Message: "registered successfully!",
		Success: true,
	}, nil
}

func (s *server) Login(ctx context.Context, req *proto.LoginRequest) (*proto.ReportResponse, error) {
	username := req.GetUsername()
	password := req.GetPassword()

	var passwordHash string
	var userId string
	err := s.db.QueryRowContext(ctx, "SELECT id, password_hash FROM users WHERE username = $1", username).Scan(&userId, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return &proto.ReportResponse{
				Message: "User not found",
				Success: false,
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	if passwordHash != password {
		return &proto.ReportResponse{
			Message: "Invalid password",
			Success: false,
		}, nil
	}

	_, _ = s.db.ExecContext(ctx, "UPDATE users SET is_online = true WHERE username = $1", username)

	sessionToken := uuid.New().String()
	err = s.rdb.Set(ctx, sessionToken, userId, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}

	return &proto.ReportResponse{
		Message: "Login successful",
		Success: true,
		Id:      sessionToken,
	}, nil

}

func (s *server) ShowUserList(ctx context.Context, req *proto.Empty) (*proto.UserListResponse, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT username, is_online FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user list: %v", err)
	}

	var userList []string
	var statusList []bool

	for rows.Next() {
		var username string
		var isOnline bool
		if err := rows.Scan(&username, &isOnline); err != nil {
			continue
		}

		userList = append(userList, username)
		statusList = append(statusList, isOnline)
	}

	return &proto.UserListResponse{
		UserList: userList,
		Status:   statusList,
	}, nil
}

func (s *server) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.ReportResponse, error) {
	sessionToken := req.GetSenderId()
	receiverUsername := req.GetReceiverId()
	message := req.GetMessage()

	senderId, err := s.rdb.Get(ctx, sessionToken).Result()
	if err != nil {
		return &proto.ReportResponse{Message: "Unauthorized: Invalid or expired session", Success: false}, nil
	}

	var ReceiverId string
	err = s.db.QueryRowContext(ctx, "SELECT id FROM users WHERE username = $1", receiverUsername).Scan(&ReceiverId)

	if err != nil {
		if err == sql.ErrNoRows {
			return &proto.ReportResponse{
				Message: "Error: Receiver username not found",
				Success: false,
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	_, err = s.db.ExecContext(ctx, "INSERT INTO messages (sender_id, receiver_id, content) VALUES ($1, $2, $3)", senderId, ReceiverId, message)
	if err != nil {
		log.Printf("Warning: Failed to send message: %v", err)
		return &proto.ReportResponse{
			Message: "Failed to send message.",
			Success: false,
		}, nil
	}

	return &proto.ReportResponse{
		Message: "Message sent!",
		Success: true,
	}, nil
}

func (s *server) Logout(ctx context.Context, req *proto.LogoutRequest) (*proto.ReportResponse, error) {
	sessionToken := req.GetId()

	userId, err := s.rdb.Get(ctx, sessionToken).Result()
	if err == nil {
		_, _ = s.db.ExecContext(ctx, "UPDATE users SET is_online = false WHERE id = $1", userId)
		s.rdb.Del(ctx, sessionToken)
	}
	return &proto.ReportResponse{
		Message: "Logged out successful",
		Success: true,
	}, nil
}

func (s *server) ShowMessage(ctx context.Context, req *proto.ShowMessagesRequest) (*proto.ShowMessagesResponse, error) {
	count := req.GetCount()
	sessionToken := req.GetId()
	contactUsername := req.GetContactId()

	userId, err := s.rdb.Get(ctx, sessionToken).Result()
	if err != nil {
		return &proto.ShowMessagesResponse{
			Message: []string{"Unauthorized: Invalid or expired session"},
			Success: false,
		}, nil
	}

	var actualContactId string
	err = s.db.QueryRowContext(ctx, "SELECT id FROM users WHERE username = $1", contactUsername).Scan(&actualContactId)

	if err != nil {
		if err == sql.ErrNoRows {
			return &proto.ShowMessagesResponse{
				Message: []string{"Error: Contact username not found"},
				Success: false,
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	query := `
		SELECT content, sender_id, seen, id FROM messages 
		WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1)
		ORDER BY id DESC 
		LIMIT $3
	`
	rows, err := s.db.QueryContext(ctx, query, userId, actualContactId, count)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch messages: %v", err)
	}
	defer rows.Close()

	var messages []string
	for rows.Next() {
		var msg string
		var sender string
		var seen bool
		var msgId int
		if err := rows.Scan(&msg, &sender, &seen, &msgId); err != nil {
			log.Printf("Warning: Failed to scan message: %v", err)
			continue
		}
		if sender != userId {
			msg = "+ " + msg
			if !seen {
				_, _ = s.db.ExecContext(ctx, "UPDATE messages SET seen = true WHERE id = $1", msgId)
			}
		} else {
			msg = "- " + msg
			if seen {
				msg = msg + "   (seen)"
			}
		}

		messages = append([]string{msg}, messages...)
	}

	return &proto.ShowMessagesResponse{
		Message: messages,
		Success: true,
	}, nil
}
