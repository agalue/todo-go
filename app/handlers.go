package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"todo-api/app/models"
)

func (a *App) indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to TODO API")
}

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

func (a *App) getTodosHandler(w http.ResponseWriter, r *http.Request) {
	if todos, err := a.db.GetAll(); err == nil {
		sendJSON(w, todos)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) getTodoHandler(w http.ResponseWriter, r *http.Request) {
	if id := getID(w, r); id > 0 {
		if todo, err := a.db.Get(id); err == nil {
			sendJSON(w, todo)
		} else {
			handleError(w, err)
		}
	}
}

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

func (a *App) deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	if id := getID(w, r); id > 0 {
		if err := a.db.Delete(id); err == nil {
			w.WriteHeader(http.StatusNoContent)
		} else {
			handleError(w, err)
		}
	}
}
