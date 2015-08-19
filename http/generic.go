package http

import (
	"github.com/Sirupsen/logrus"
	"net/http"
)

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func HandleError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	logrus.WithField("error", err).Error("error in task list handler")
}
