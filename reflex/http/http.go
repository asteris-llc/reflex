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

type Registerer interface {
	Register(*mux.Router)
}

type HTTP struct {
	Components []Registerer
}

func (h *HTTP) ServeHTTP(addr string) {
	router := mux.NewRouter()

	for _, component := range h.Components {
		component.Register(router)
	}

	n := negroni.New()

	// middleware
	n.Use(negroni.NewRecovery())
	n.Use(negronilogrus.NewMiddleware())
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(router)

	// start
	logrus.WithField("addr", addr).Info("listening")
	graceful.Run(addr, 10*time.Second, n)
}
