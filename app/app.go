package app

import (
	"context"
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"todo-api/app/database"

	_ "todo-api/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

//go:embed web/dist
var web embed.FS

type App struct {
	db     database.TodoDB
	router *http.ServeMux
	server *http.Server
}

func New(db database.TodoDB) *App {
	a := &App{
		db:     db,
		router: http.NewServeMux(),
	}
	a.initRoutes()
	return a
}

func (a *App) initRoutes() {
	a.router.HandleFunc("POST /api/v1/todos", a.addTodoHandler)
	a.router.HandleFunc("GET /api/v1/todos", a.getTodosHandler)
	a.router.HandleFunc("GET /api/v1/todos/{id}", a.getTodoHandler)
	a.router.HandleFunc("PUT /api/v1/todos/{id}", a.updateTodoHandler)
	a.router.HandleFunc("DELETE /api/v1/todos/{id}", a.deleteTodoHandler)
	a.router.Handle("GET /swagger/*", httpSwagger.Handler())

	dist, err := fs.Sub(web, "web/dist")
	if err != nil {
		slog.Error("cannot mount web interfce", slog.String("error", err.Error()))
		return
	}
	a.router.Handle("/", http.FileServer(http.FS(dist)))
}

func (a *App) Start(addr string) {
	a.server = &http.Server{
		Addr:    addr,
		Handler: NewObserver(a.router),
	}

	if err := a.db.Init(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	go func() {
		slog.Info("starting server", "address", addr)
		if err := a.server.ListenAndServe(); err != nil {
			slog.Warn(err.Error())
			if err != http.ErrServerClosed {
				os.Exit(1)
			}
		}
	}()
}

func (a *App) Shutdown() {
	if err := a.server.Shutdown(context.Background()); err != nil {
		slog.Warn(err.Error())
	}
}
