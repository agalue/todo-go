package database

import (
	"context"
	"testing"
	"todo-api/app/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gotest.tools/v3/assert"
)

func initMockDatabase() (*gorm.DB, sqlmock.Sqlmock) {
	mockDb, mockObj, _ := sqlmock.New()
	dialector := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: "postgres",
	})
	db, _ := gorm.Open(dialector, &gorm.Config{
		Logger:                 logger.Default.LogMode(logger.Silent),
		SkipDefaultTransaction: true,
	})
	return db, mockObj
}

func TestGetAll(t *testing.T) {
	cli, mock := initMockDatabase()
	rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Pass the test")
	mock.ExpectQuery(`^SELECT .*`).WillReturnRows(rows)

	db := &DB{cli: cli}
	todos, err := db.GetAll(context.Background())
	assert.NilError(t, err)
	assert.Equal(t, 1, len(todos))
}

func TestGet(t *testing.T) {
	cli, mock := initMockDatabase()
	rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Pass the test")
	mock.ExpectQuery(`^SELECT . FROM "todos" WHERE "todos"."id" .*`).WillReturnRows(rows)

	db := &DB{cli: cli}

	todo, err := db.Get(context.Background(), 1)
	assert.NilError(t, err)
	assert.Equal(t, 1, todo.ID)
}

func TestGetNotFound(t *testing.T) {
	cli, mock := initMockDatabase()
	mock.ExpectQuery(`^SELECT . FROM "todos" WHERE "todos"."id" .*`).WillReturnError(gorm.ErrRecordNotFound)

	db := &DB{cli: cli}

	_, err := db.Get(context.Background(), 2)
	assert.Equal(t, ErrorNotFound, err)
}

func TestAdd(t *testing.T) {
	cli, mock := initMockDatabase()
	mock.ExpectQuery(`INSERT .* RETURNING .*`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	db := &DB{cli: cli}

	todo, err := db.Add(context.Background(), models.Base{Title: "Testing"})
	assert.NilError(t, err)
	assert.Equal(t, 1, todo.ID)
}

func TestSetStatus(t *testing.T) {
	cli, mock := initMockDatabase()
	rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Pass the test")
	mock.ExpectQuery(`^SELECT . FROM "todos" WHERE "todos"."id" .*`).WillReturnRows(rows)
	mock.ExpectExec(`^UPDATE .*`).WillReturnResult(sqlmock.NewResult(1, 1))
	db := &DB{cli: cli}

	err := db.SetStatus(context.Background(), 1, models.Status{Completed: true})
	assert.NilError(t, err)
}

func TestSetStatusNotFound(t *testing.T) {
	cli, mock := initMockDatabase()
	mock.ExpectQuery(`^SELECT . FROM "todos" WHERE "todos"."id" .*`).WillReturnError(gorm.ErrRecordNotFound)
	db := &DB{cli: cli}

	err := db.SetStatus(context.Background(), 1, models.Status{Completed: true})
	assert.Equal(t, ErrorNotFound, err)
}

func TestDelete(t *testing.T) {
	cli, mock := initMockDatabase()
	rows := sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Pass the test")
	mock.ExpectQuery(`^SELECT . FROM "todos" WHERE "todos"."id" .*`).WillReturnRows(rows)
	mock.ExpectExec(`^DELETE .*`).WillReturnResult(sqlmock.NewResult(1, 1))
	db := &DB{cli: cli}

	err := db.Delete(context.Background(), 1)
	assert.NilError(t, err)
}

func TestDeleteNotFound(t *testing.T) {
	cli, mock := initMockDatabase()
	mock.ExpectQuery(`^SELECT . FROM "todos" WHERE "todos"."id" .*`).WillReturnError(gorm.ErrRecordNotFound)
	db := &DB{cli: cli}

	err := db.Delete(context.Background(), 1)
	assert.Equal(t, ErrorNotFound, err)
}
