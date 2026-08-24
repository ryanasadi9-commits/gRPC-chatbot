package main

import (
	"hamrahTask1/internal/app"
	"hamrahTask1/pkg/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	cfg := config.LoadEnv()

	application := app.New(cfg)

	application.Start()
}
