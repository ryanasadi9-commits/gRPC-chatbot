package ports

import "context"

type LogEntry struct {
	Level       string
	Message     string
	ErrorDetail string
}

type LogRepository interface {
	SaveLogs(ctx context.Context, entries []LogEntry) error
}
