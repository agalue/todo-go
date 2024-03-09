package database

import (
	"context"
	"errors"
	"todo-api/app/models"
)

var ErrorNotFound = errors.New("record not found")

type TodoDB interface {
	Init() error
	Shutdown()
	GetAll(ctx context.Context) ([]models.Todo, error)
	Get(ctx context.Context, id int) (models.Todo, error)
	Add(ctx context.Context, todo models.Base) (models.Todo, error)
	SetStatus(ctx context.Context, id int, status models.Status) error
	Delete(ctx context.Context, id int) error
}
