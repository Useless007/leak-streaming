package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pressly/goose/v3"

	"github.com/leak-streaming/leak-streaming/backend/internal/platform/config"
	"github.com/leak-streaming/leak-streaming/backend/internal/platform/database"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	defer db.Close()

	command := "up"
	args := os.Args[1:]
	if len(args) > 0 {
		command = args[0]
		args = args[1:]
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set dialect: %v", err)
	}

	if err := goose.RunContext(ctx, command, db, cfg.Database.MigrationsDir, args...); err != nil {
		log.Fatalf("migrate command %q failed: %v", command, err)
	}

	fmt.Printf("goose %s succeeded\n", command)
}
