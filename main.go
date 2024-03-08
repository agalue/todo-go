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

// docker run -d --rm --name postgres -e POSTGRES_DB=todo -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 postgres:16-alpine

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := getEnv("API_LISTEN", ":8080")
	dsn := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		getEnv("PGHOST", "localhost"),
		getEnv("PGPORT", "5432"),
		getEnv("PGDATABASE", "todo"),
		getEnv("PGUSER", "postgres"),
		getEnv("PGPASSWORD", "postgres"))

	db := database.New(dsn)
	defer db.Shutdown()

	a := app.New(db)
	defer a.Shutdown()

	a.Start(addr)
	<-ctx.Done()
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
