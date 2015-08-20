package http

import (
	"github.com/asteris-llc/reflex/reflex/logic"
	"github.com/gorilla/mux"
)

type API struct {
	Events chan *logic.Event
}

func NewAPI(events chan *logic.Event) (*API, error) {
	return &API{events}, nil
}

func (a *API) Register(router *mux.Router) {
	v1 := router.PathPrefix("/1").Subrouter()

	// tasks
	// TODO: task registration

	// events
	eventsHandler := &EventsHandler{a.Events}
	events := v1.PathPrefix("/events").Subrouter()
	events.Methods("POST").HandlerFunc(eventsHandler.Create)
	events.Methods("GET", "HEAD", "PUT", "PATCH", "DELETE").HandlerFunc(MethodNotAllowed)
}
