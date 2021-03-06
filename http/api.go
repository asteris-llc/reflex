package http

import (
	"github.com/Sirupsen/logrus"
	"github.com/asteris-llc/reflex/state"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/tylerb/graceful"
	"time"
)

type API struct {
	store  state.Storer
	Events chan *state.Event
}

func NewAPI(store state.Storer) *API {
	return &API{store, make(chan *state.Event)}
}

func (a *API) Serve(addr string) {
	router := mux.NewRouter()

	v1 := router.PathPrefix("/1").Subrouter()

	// tasks
	tasksHandler := TasksHandler{a.store.Tasks()}

	task := v1.PathPrefix("/tasks/{id}").Subrouter()
	task.Methods("GET").HandlerFunc(tasksHandler.Get)
	task.Methods("PUT").HandlerFunc(tasksHandler.Set)
	task.Methods("DELETE").HandlerFunc(tasksHandler.Delete)
	task.Methods("POST", "PATCH").HandlerFunc(MethodNotAllowed)

	tasks := v1.PathPrefix("/tasks").Subrouter()
	tasks.Methods("GET").HandlerFunc(tasksHandler.List)
	tasks.Methods("POST").HandlerFunc(tasksHandler.Set)
	tasks.Methods("PUT", "DELETE", "PATCH").HandlerFunc(MethodNotAllowed)

	// events
	eventsHandler := EventsHandler{a.store.Events(), a.Events}

	event := v1.PathPrefix("/events/{id}").Subrouter()
	event.Methods("GET").HandlerFunc(eventsHandler.Get)
	event.Methods("PUT", "POST", "DELETE", "PATCH").HandlerFunc(MethodNotAllowed)

	events := v1.PathPrefix("/events").Subrouter()
	events.Methods("GET").HandlerFunc(eventsHandler.List)
	events.Methods("POST").HandlerFunc(eventsHandler.Create)
	events.Methods("PUT", "DELETE", "PATCH").HandlerFunc(MethodNotAllowed)

	// set up event handlers

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negronilogrus.NewMiddleware())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(router)

	logrus.WithField("addr", addr).Info("listening")
	graceful.Run(addr, 10*time.Second, n)
}
