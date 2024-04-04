package main

import (
	"context"
	"os/signal"
	"syscall"
	"todo-api/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a := app.New()
	a.Start(ctx)
	<-ctx.Done()
	a.Shutdown()
}
