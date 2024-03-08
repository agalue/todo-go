package app

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-api/app/database"
)

func getID(w http.ResponseWriter, r *http.Request) int {
	if id, err := strconv.Atoi(r.PathValue("id")); err == nil {
		return id
	} else {
		http.Error(w, "Invalid ID", http.StatusForbidden)
	}
	return -1
}

func sendJSON(w http.ResponseWriter, obj any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(obj)
}

func handleError(w http.ResponseWriter, err error) {
	if err == database.ErrorNotFound {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
