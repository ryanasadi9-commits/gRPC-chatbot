package ports

import "context"

type ChatRepository interface {
	GetIDByUsername(ctx context.Context, username string) (string, error)
	SaveMessage(ctx context.Context, senderID, receiverID, content string) error
	GetMessages(ctx context.Context, userID, contactID string, count int64) ([]string, error)
}
