package http

import (
	"github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	negronilogrus "github.com/meatballhat/negroni-logrus"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/tylerb/graceful"
	"time"
)

type API struct{}

func (a *API) Serve(addr string) {
	router := mux.NewRouter()

	v1 := router.PathPrefix("/1").Subrouter()

	tasksHandler := TasksHandler{}
	tasks := v1.PathPrefix("/tasks").Subrouter()
	tasks.Methods("GET").HandlerFunc(tasksHandler.List)
	tasks.Methods("PUT", "POST", "DELETE", "PATCH").HandlerFunc(MethodNotAllowed)

	task := tasks.PathPrefix("/{id}").Subrouter()
	task.Methods("GET").HandlerFunc(tasksHandler.Get)
	task.Methods("PUT").HandlerFunc(tasksHandler.Set)
	task.Methods("DELETE").HandlerFunc(tasksHandler.Delete)
	task.Methods("POST", "PATCH").HandlerFunc(MethodNotAllowed)

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negronilogrus.NewMiddleware())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(router)

	logrus.WithField("addr", addr).Info("listening")
	graceful.Run(addr, 10*time.Second, n)
}
