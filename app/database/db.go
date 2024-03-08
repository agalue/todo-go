package database

import (
	"errors"
	"todo-api/app/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var ErrorNotFound = errors.New("record not found")

type TodoDB interface {
	Init() error
	Shutdown()
	GetAll() ([]models.Todo, error)
	Get(id int) (models.Todo, error)
	Add(todo models.Base) (models.Todo, error)
	SetStatus(id int, status models.Status) error
	Delete(id int) error
}

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

func (db *DB) GetAll() ([]models.Todo, error) {
	todos := make([]models.Todo, 0)
	err := db.cli.Find(&todos).Error
	return todos, err
}

func (db *DB) Get(id int) (models.Todo, error) {
	todo := models.Todo{}
	err := db.cli.First(&todo, id).Error
	if err == gorm.ErrRecordNotFound {
		err = ErrorNotFound
	}
	return todo, err
}

func (db *DB) Add(todo models.Base) (models.Todo, error) {
	dbtodo := models.Todo{Base: todo}
	err := db.cli.Create(&dbtodo).Error
	return dbtodo, err
}

func (db *DB) SetStatus(id int, status models.Status) error {
	if todo, err := db.Get(id); err == nil {
		todo.Completed = status.Completed
		return db.cli.Save(todo).Error
	} else {
		return err
	}
}

func (db *DB) Delete(id int) error {
	if todo, err := db.Get(id); err == nil {
		return db.cli.Delete(todo).Error
	} else {
		return err
	}
}
