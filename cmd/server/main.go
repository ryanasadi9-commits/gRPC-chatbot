package main

import (
	"hamrahTask1/internal/app"
	"hamrahTask1/pkg/config"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.LoadEnv()

	application := app.New(cfg)

	go application.Start()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	application.Stop()
}
