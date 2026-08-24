package postgress

import (
	"context"
	"database/sql"
	"hamrahTask1/internal/core/ports"
)

type PostgresLogRepository struct {
	db *sql.DB
}

func NewPostgresLogRepository(db *sql.DB) ports.LogRepository {
	return &PostgresLogRepository{db: db}
}

func (r *PostgresLogRepository) SaveLogs(ctx context.Context, entries []ports.LogEntry) error {
	if len(entries) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, "INSERT INTO app_logs (level, message, error_detail) VALUES ($1, $2, $3)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, entry := range entries {
		_, err = stmt.ExecContext(ctx, entry.Level, entry.Message, entry.ErrorDetail)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
