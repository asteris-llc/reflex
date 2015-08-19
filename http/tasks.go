package http

import (
	"encoding/json"
	"github.com/asteris-llc/reflex/state"
	"github.com/gorilla/mux"
	"net/http"
)

type TasksHandler struct {
	store state.TaskStorer
}

func (t *TasksHandler) List(w http.ResponseWriter, r *http.Request) {
	tasks, err := t.store.List()
	if err != nil {
		HandleError(err, w)
		return
	}

	blob, err := json.Marshal(tasks)
	if err != nil {
		HandleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(blob)
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
