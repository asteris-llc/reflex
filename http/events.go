package http

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"github.com/asteris-llc/reflex/state"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type EventsHandler struct {
	events    state.EventStorer
	NewEvents chan *state.Event
}

func (e *EventsHandler) List(w http.ResponseWriter, r *http.Request) {
	events, err := e.events.List()
	if err != nil {
		HandleError(err, w)
		return
	}

	// TODO: probably paginate this
	blob, err := json.Marshal(events)
	if err != nil {
		HandleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(blob)
}

func (e *EventsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	event, err := e.events.Get(id)
	if err == state.ErrNoEvent {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	} else if err != nil {
		HandleError(err, w)
		return
	}

	blob, err := json.Marshal(event)
	if err != nil {
		HandleError(err, w)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(blob)
}

func (e *EventsHandler) Create(w http.ResponseWriter, r *http.Request) {
	blob, err := ioutil.ReadAll(r.Body)
	if err != nil {
		HandleError(err, w)
		return
	}
	defer r.Body.Close()

	event := new(state.Event)
	err = json.Unmarshal(blob, event)
	if err != nil {
		HandleError(err, w)
		return
	}

	event.ID = uuid.NewRandom().String()

	err = e.events.Update(event)
	if err != nil {
		HandleError(err, w)
		return
	}

	e.NewEvents <- event

	headers := w.Header()
	headers.Add("Location", "/1/events/"+event.ID)
	w.WriteHeader(http.StatusCreated)
}
