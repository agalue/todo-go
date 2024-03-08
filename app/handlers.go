package app

import (
	"encoding/json"
	"net/http"
	"todo-api/app/models"
)

// @Summary Add a new TODO
// @Consume json
// @Produce json
// @Param   todo body models.Base true "New TODO"
// @Success 201 {object} models.Todo
// @Failure 500 "Backend error"
// @Router  /api/v1/todos [post]
func (a *App) addTodoHandler(w http.ResponseWriter, r *http.Request) {
	base := models.Base{}
	if err := json.NewDecoder(r.Body).Decode(&base); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if todo, err := a.db.Add(base); err == nil {
		w.WriteHeader(http.StatusCreated)
		sendJSON(w, todo)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Get all the TODOs
// @Produce json
// @Success 200 {object} []models.Todo
// @Failure 500 "Backend error"
// @Router  /api/v1/todos [get]
func (a *App) getTodosHandler(w http.ResponseWriter, r *http.Request) {
	if todos, err := a.db.GetAll(); err == nil {
		sendJSON(w, todos)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// @Summary Get a TODO
// @Produce json
// @Param   id path int true "TODO ID"
// @Success 200 {object} models.Todo
// @Failure 404 "Not found"
// @Failure 500 "Backend error"
// @Router  /api/v1/todos/{id} [get]
func (a *App) getTodoHandler(w http.ResponseWriter, r *http.Request) {
	if id := getID(w, r); id > 0 {
		if todo, err := a.db.Get(id); err == nil {
			sendJSON(w, todo)
		} else {
			handleError(w, err)
		}
	}
}

// @Summary Update a TODO
// @Produce json
// @Param   id path int true "TODO ID"
// @Param   status body models.Status true "TODO Status"
// @Success 200 {object} models.Todo
// @Failure 404 "Not found"
// @Failure 500 "Backend error"
// @Router  /api/v1/todos/{id} [put]
func (a *App) updateTodoHandler(w http.ResponseWriter, r *http.Request) {
	if id := getID(w, r); id > 0 {
		status := models.Status{}
		if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := a.db.SetStatus(id, status); err == nil {
			w.WriteHeader(http.StatusNoContent)
		} else {
			handleError(w, err)
		}
	}
}

// @Summary Delete a TODO
// @Produce json
// @Param   id path int true "TODO ID"
// @Success 200 {object} models.Todo
// @Failure 404 "Not found"
// @Failure 500 "Backend error"
// @Router  /api/v1/todos/{id} [delete]
func (a *App) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	if id := getID(w, r); id > 0 {
		if err := a.db.Delete(id); err == nil {
			w.WriteHeader(http.StatusNoContent)
		} else {
			handleError(w, err)
		}
	}
}
