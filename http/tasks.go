package http

import (
	"encoding/json"
	"github.com/asteris-llc/reflex/state"
	"github.com/gorilla/mux"
	"io/ioutil"
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

	task, err := t.store.Get(id)
	if err == state.ErrNoTask {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	} else if err != nil {
		HandleError(err, w)
		return
	}

	blob, err := json.Marshal(task)
	if err != nil {
		HandleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(blob)
}

func (t *TasksHandler) Set(w http.ResponseWriter, r *http.Request) {
	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HandleError(err, w)
		return
	}
	defer r.Body.Close()

	task := new(state.Task)
	err = json.Unmarshal(blob, task)
	if err != nil {
		HandleError(err, w)
		return
	}

	// validate
	if id, ok := mux.Vars(r)["id"]; ok && id != task.ID {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("body ID does not match URL ID"))
		return
	}

	err = t.store.Update(task)
	if err != nil {
		HandleError(err, w)
		return
	}

	headers := w.Header()
	headers.Add("Location", "/1/tasks/"+task.ID)
	w.WriteHeader(http.StatusCreated)
}

func (t *TasksHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := t.store.Delete(id)
	if err != nil {
		HandleError(err, w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
