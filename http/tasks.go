package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

type TasksHandler struct{}

func (t *TasksHandler) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (t *TasksHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

func (t *TasksHandler) Set(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

func (t *TasksHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// id := mux.Vars(r)["id"]

	w.WriteHeader(http.StatusNoContent)
}
