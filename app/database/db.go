package database

import (
	"context"
	"todo-api/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DB struct {
	dsn string
	cli *gorm.DB
}

func New(dsn string) TodoDB {
	return &DB{
		dsn: dsn,
	}
}

func (db *DB) Init() error {
	var err error
	if db.cli, err = gorm.Open(postgres.Open(db.dsn), &gorm.Config{}); err != nil {
		return err
	}
	if err = db.cli.AutoMigrate(&models.Todo{}); err != nil {
		return err
	}
	return nil
}

func (db *DB) Shutdown() {
}

func (db *DB) GetAll(ctx context.Context) ([]models.Todo, error) {
	todos := make([]models.Todo, 0)
	err := db.cli.WithContext(ctx).Find(&todos).Error
	return todos, err
}

func (db *DB) Get(ctx context.Context, id int) (models.Todo, error) {
	todo := models.Todo{}
	err := db.cli.WithContext(ctx).First(&todo, id).Error
	if err == gorm.ErrRecordNotFound {
		err = ErrorNotFound
	}
	return todo, err
}

func (db *DB) Add(ctx context.Context, todo models.Base) (models.Todo, error) {
	dbtodo := models.Todo{Base: todo}
	err := db.cli.WithContext(ctx).Create(&dbtodo).Error
	return dbtodo, err
}

func (db *DB) SetStatus(ctx context.Context, id int, status models.Status) error {
	if todo, err := db.Get(ctx, id); err == nil {
		todo.Completed = status.Completed
		return db.cli.WithContext(ctx).Save(todo).Error
	} else {
		return err
	}
}

func (db *DB) Delete(ctx context.Context, id int) error {
	if todo, err := db.Get(ctx, id); err == nil {
		return db.cli.WithContext(ctx).Delete(todo).Error
	} else {
		return err
	}
}
