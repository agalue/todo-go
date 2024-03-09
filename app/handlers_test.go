package app

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"todo-api/app/database"
	"todo-api/app/models"

	"gotest.tools/v3/assert"
)

var ErrorMockInternal = errors.New("something went wrong")

type MockDB struct {
	todos []*models.Todo
	fail  bool
}

func (db *MockDB) Init() error {
	db.todos = []*models.Todo{
		{
			ID: 0,
			Base: models.Base{
				Title: "Test API",
			},
		},
		{
			ID: 1,
			Base: models.Base{
				Title: "Test DB",
			},
		},
	}
	return nil
}

func (db *MockDB) Shutdown() {
}

func (db *MockDB) GetAll(ctx context.Context) ([]models.Todo, error) {
	if db.fail {
		return nil, ErrorMockInternal
	}
	todos := make([]models.Todo, len(db.todos))
	for i, todo := range db.todos {
		todos[i] = *todo
	}
	return todos, nil
}

func (db *MockDB) Get(ctx context.Context, id int) (models.Todo, error) {
	if db.fail {
		return models.Todo{}, ErrorMockInternal
	}
	if id < len(db.todos) {
		return *db.todos[id], nil
	}
	return models.Todo{}, database.ErrorNotFound
}

func (db *MockDB) Add(ctx context.Context, base models.Base) (models.Todo, error) {
	if db.fail {
		return models.Todo{}, ErrorMockInternal
	}
	todo := models.Todo{Base: base}
	db.todos = append(db.todos, &todo)
	todo.ID = len(db.todos) - 1
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	return todo, nil
}

func (db *MockDB) SetStatus(ctx context.Context, id int, status models.Status) error {
	if db.fail {
		return ErrorMockInternal
	}
	if id < len(db.todos) {
		db.todos[id].Completed = status.Completed
		return nil
	}
	return database.ErrorNotFound
}

func (db *MockDB) Delete(ctx context.Context, id int) error {
	if db.fail {
		return ErrorMockInternal
	}
	if _, err := db.Get(ctx, id); err == nil {
		db.todos = append(db.todos[:id], db.todos[id+1:]...)
		return nil
	}
	return database.ErrorNotFound
}

func TestAddTodoHandler(t *testing.T) {
	db := new(MockDB)
	db.Init()
	srv := New(db)

	base := models.Base{
		Title: "TestAddTodoHandler",
	}
	data, err := json.Marshal(base)
	assert.NilError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/todos", bytes.NewBuffer(data))
	w := httptest.NewRecorder()

	srv.addTodoHandler(w, r)
	resp := w.Result()
	todo := new(models.Todo)
	err = json.NewDecoder(resp.Body).Decode(todo)
	assert.NilError(t, err)

	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, 2, todo.ID)
	assert.Equal(t, 3, len(db.todos))
}

func TestAddTodoHandlerInvalidData(t *testing.T) {
	db := new(MockDB)
	db.Init()
	srv := New(db)

	r := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer([]byte("Invalid Data")))
	w := httptest.NewRecorder()

	srv.addTodoHandler(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAddTodoHandlerServerError(t *testing.T) {
	db := &MockDB{fail: true}
	db.Init()
	srv := New(db)

	base := models.Base{
		Title: "TestAddTodoHandler",
	}
	data, err := json.Marshal(base)
	assert.NilError(t, err)

	r := httptest.NewRequest(http.MethodPost, "/api/v1/todos", bytes.NewBuffer(data))
	w := httptest.NewRecorder()

	srv.addTodoHandler(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetTodosHandler(t *testing.T) {
	db := new(MockDB)
	db.Init()
	srv := New(db)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	w := httptest.NewRecorder()

	srv.getTodosHandler(w, r)

	resp := w.Result()
	todos := make([]models.Todo, 0)
	err := json.NewDecoder(resp.Body).Decode(&todos)
	assert.NilError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 2, len(todos))
}

func TestGetTodosHandlerServerError(t *testing.T) {
	db := &MockDB{fail: true}
	db.Init()
	srv := New(db)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/todos", nil)
	w := httptest.NewRecorder()

	srv.getTodosHandler(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetTodoHandler(t *testing.T) {
	db := new(MockDB)
	db.Init()
	srv := New(db)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/todos/1", nil)
	w := httptest.NewRecorder()
	r.SetPathValue("id", "1")

	srv.getTodoHandler(w, r)

	resp := w.Result()
	todo := new(models.Todo)
	err := json.NewDecoder(resp.Body).Decode(todo)
	assert.NilError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, 1, todo.ID)
}

func TestGetTodoHandlerNotFound(t *testing.T) {
	db := new(MockDB)
	db.Init()
	srv := New(db)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/todos/10", nil)
	w := httptest.NewRecorder()
	r.SetPathValue("id", "10")

	srv.getTodoHandler(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestGetTodoHandlerServerError(t *testing.T) {
	db := &MockDB{fail: true}
	db.Init()
	srv := New(db)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/todos/1", nil)
	w := httptest.NewRecorder()
	r.SetPathValue("id", "1")

	srv.getTodoHandler(w, r)
	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestUpdateTodosHandler(t *testing.T) {
	db := new(MockDB)
	db.Init()
	srv := New(db)

	base := models.Status{
		Completed: true,
	}
	data, err := json.Marshal(base)
	assert.NilError(t, err)

	r := httptest.NewRequest(http.MethodPut, "/api/v1/todos/1", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	r.SetPathValue("id", "1")

	srv.updateTodoHandler(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Assert(t, db.todos[1].Completed)
}

func TestUpdateTodosHandlerInvalidData(t *testing.T) {
	db := new(MockDB)
	db.Init()
	srv := New(db)

	r := httptest.NewRequest(http.MethodPut, "/api/v1/todos/1", bytes.NewBuffer([]byte("Invalid Data")))
	w := httptest.NewRecorder()
	r.SetPathValue("id", "1")

	srv.updateTodoHandler(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	assert.Assert(t, db.todos[1].Completed == false)
}

func TestUpdateTodosHandlerNotFound(t *testing.T) {
	db := new(MockDB)
	db.Init()
	srv := New(db)

	base := models.Status{
		Completed: true,
	}
	data, err := json.Marshal(base)
	assert.NilError(t, err)

	r := httptest.NewRequest(http.MethodGet, "/api/v1/todos/10", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	r.SetPathValue("id", "10")

	srv.updateTodoHandler(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestUpdateTodosHandlerServerError(t *testing.T) {
	db := &MockDB{fail: true}
	db.Init()
	srv := New(db)

	base := models.Status{
		Completed: true,
	}
	data, err := json.Marshal(base)
	assert.NilError(t, err)

	r := httptest.NewRequest(http.MethodPut, "/api/v1/todos/1", bytes.NewBuffer(data))
	w := httptest.NewRecorder()
	r.SetPathValue("id", "1")

	srv.updateTodoHandler(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Assert(t, db.todos[1].Completed == false)
}

func TestDeleteTodosHandler(t *testing.T) {
	db := new(MockDB)
	db.Init()
	srv := New(db)

	r := httptest.NewRequest(http.MethodDelete, "/api/v1/todos/1", nil)
	w := httptest.NewRecorder()
	r.SetPathValue("id", "1")

	srv.deleteTodoHandler(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	assert.Equal(t, 1, len(db.todos))
}

func TestDeleteTodosHandlerNotFound(t *testing.T) {
	db := new(MockDB)
	db.Init()
	srv := New(db)

	r := httptest.NewRequest(http.MethodDelete, "/api/v1/todos/10", nil)
	w := httptest.NewRecorder()
	r.SetPathValue("id", "10")

	srv.deleteTodoHandler(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestDeleteTodosHandlerServerError(t *testing.T) {
	db := &MockDB{fail: true}
	db.Init()
	srv := New(db)

	r := httptest.NewRequest(http.MethodDelete, "/api/v1/todos/1", nil)
	w := httptest.NewRecorder()
	r.SetPathValue("id", "1")

	srv.deleteTodoHandler(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	assert.Equal(t, 2, len(db.todos))
}
