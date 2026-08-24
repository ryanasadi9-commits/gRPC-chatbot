package ports

import "context"

type UserRepository interface {
	CreateUser(ctx context.Context, username, password string, age int64, gender bool) error
	GetUserByUsername(ctx context.Context, username string) (id string, hash string, err error)
	UpdateOnlineStatus(ctx context.Context, identifier string, isOnline bool, byUsername bool) error
	GetAllUsers(ctx context.Context) ([]string, []bool, error)
}

type SessionCache interface {
	CreateSession(ctx context.Context, token, userID string) error
	GetUserID(ctx context.Context, token string) (string, error)
	DeleteSession(ctx context.Context, token string) error
}
