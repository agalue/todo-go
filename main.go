package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"todo-api/app"
	"todo-api/app/database"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := getEnv("API_LISTEN", ":8080")
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		getEnv("POSTGRES_HOST", "localhost"),
		getEnv("POSTGRES_PORT", "5432"),
		getEnv("POSTGRES_DB", "todo"),
		getEnv("POSTGRES_USER", "postgres"),
		getEnv("POSTGRES_PASSWORD", "postgres"))

	db := database.New(dsn)
	a := app.New(db)
	a.Start(ctx, addr)
	<-ctx.Done()
	a.Shutdown()
	db.Shutdown()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
