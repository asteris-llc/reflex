package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

type EventsHandler struct{}

func (t *EventsHandler) List(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func (t *EventsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

func (t *EventsHandler) Create(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
