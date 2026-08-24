package ports

import "context"

type Service interface {
	Register(ctx context.Context, username, password string, age int64, gender bool) error
	Login(ctx context.Context, username, password string) (string, error)
	Logout(ctx context.Context, token string) error
	GetUsers(ctx context.Context) ([]string, []bool, error)
	SendMessage(ctx context.Context, token, receiverUsername, message string) error
	GetMessages(ctx context.Context, token, contactUsername string, count int64) ([]string, error)
}
