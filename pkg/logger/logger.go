package logger

import (
	"context"
	"log"
	"time"

	"hamrahTask1/internal/core/ports"
)

type Logger struct {
	repo    ports.LogRepository
	logChan chan ports.LogEntry
}

func New(repo ports.LogRepository) *Logger {
	l := &Logger{
		repo:    repo,
		logChan: make(chan ports.LogEntry, 5000),
	}

	const numWorkers = 3
	for i := 0; i < numWorkers; i++ {
		go l.startBatchWorker()
	}

	return l
}

func (l *Logger) startBatchWorker() {
	var batch []ports.LogEntry
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case entry := <-l.logChan:
			batch = append(batch, entry)

			if len(batch) >= 500 {
				_ = l.repo.SaveLogs(context.Background(), batch)
				batch = nil
			}
		case <-ticker.C:
			if len(batch) > 0 {
				_ = l.repo.SaveLogs(context.Background(), batch)
				batch = nil
			}
		}
	}
}

func (l *Logger) Info(msg string) {
	log.Println("INFO: " + msg)

	l.logChan <- ports.LogEntry{Level: "INFO", Message: msg, ErrorDetail: ""}
}

func (l *Logger) Error(msg string, err error) {
	log.Printf("ERROR: %s: %v\n", msg, err)

	l.logChan <- ports.LogEntry{Level: "ERROR", Message: msg, ErrorDetail: err.Error()}
}

func (l *Logger) Fatal(msg string, err error) {
	log.Printf("FATAL: %s: %v\n", msg, err)

	_ = l.repo.SaveLogs(context.Background(), []ports.LogEntry{
		{Level: "FATAL", Message: msg, ErrorDetail: err.Error()},
	})

	log.Fatalf("FATAL: %s: %v", msg, err)
}
