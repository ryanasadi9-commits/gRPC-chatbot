package postgress

import (
	"context"
	"database/sql"
	"hamrahTask1/internal/core/ports"
)

type PostgresChatRepository struct {
	db *sql.DB
}

func NewPostgresChatRepository(db *sql.DB) ports.ChatRepository {
	return &PostgresChatRepository{db: db}
}

func (r *PostgresChatRepository) GetIDByUsername(ctx context.Context, username string) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, "SELECT id FROM users WHERE username = $1", username).Scan(&id)
	return id, err
}

func (r *PostgresChatRepository) SaveMessage(ctx context.Context, senderID, receiverID, content string) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO messages (sender_id, receiver_id, content) VALUES ($1, $2, $3)", senderID, receiverID, content)
	return err
}

func (r *PostgresChatRepository) GetMessages(ctx context.Context, userID, contactID string, count int64) ([]string, error) {
	query := `SELECT content, sender_id, seen, id FROM messages WHERE (sender_id = $1 AND receiver_id = $2) OR (sender_id = $2 AND receiver_id = $1) ORDER BY id DESC LIMIT $3`
	rows, err := r.db.QueryContext(ctx, query, userID, contactID, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []string
	for rows.Next() {
		var msg, sender string
		var seen bool
		var msgId int
		if err := rows.Scan(&msg, &sender, &seen, &msgId); err != nil {
			continue
		}
		if sender != userID {
			msg = "+ " + msg
			if !seen {
				r.db.ExecContext(ctx, "UPDATE messages SET seen = true WHERE id = $1", msgId)
			}
		} else {
			msg = "- " + msg
			if seen {
				msg += "   (seen)"
			}
		}
		msgs = append([]string{msg}, msgs...)
	}
	return msgs, nil
}
