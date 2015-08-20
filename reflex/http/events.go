package http

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/asteris-llc/reflex/reflex/logic"
	"io/ioutil"
	"net/http"
)

type EventsHandler struct {
	Events chan *logic.Event
}

func (e *EventsHandler) Create(w http.ResponseWriter, r *http.Request) {
	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HandleError(err, w)
		return
	}
	defer r.Body.Close()

	event := new(logic.Event)
	err = json.Unmarshal(blob, event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error() + "\n"))
		return
	}

	event.ID = uuid.NewRandom().String()

	if err != nil {
		HandleError(err, w)
		return
	}

	e.Events <- event

	w.WriteHeader(http.StatusCreated)
}
