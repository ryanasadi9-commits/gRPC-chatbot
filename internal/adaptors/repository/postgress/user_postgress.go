package postgress

import (
	"context"
	"database/sql"
	"hamrahTask1/internal/core/ports"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) ports.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, username, password string, age int64, gender bool) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (username, password_hash, age, gender, is_online) VALUES ($1, $2, $3, $4, $5)", username, password, age, gender, false)
	return err
}

func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (string, string, error) {
	var id, hash string
	err := r.db.QueryRowContext(ctx, "SELECT id, password_hash FROM users WHERE username = $1", username).Scan(&id, &hash)
	return id, hash, err
}

func (r *PostgresUserRepository) UpdateOnlineStatus(ctx context.Context, identifier string, isOnline bool, byUsername bool) error {
	query := "UPDATE users SET is_online = $1 WHERE id = $2"
	if byUsername {
		query = "UPDATE users SET is_online = $1 WHERE username = $2"
	}
	_, err := r.db.ExecContext(ctx, query, isOnline, identifier)
	return err
}

func (r *PostgresUserRepository) GetAllUsers(ctx context.Context) ([]string, []bool, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT username, is_online FROM users")
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	var users []string
	var statuses []bool
	for rows.Next() {
		var u string
		var s bool
		if err := rows.Scan(&u, &s); err == nil {
			users = append(users, u)
			statuses = append(statuses, s)
		}
	}
	return users, statuses, nil
}
