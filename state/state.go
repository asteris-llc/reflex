package state

import (
	"errors"
)

var (
	ErrNoEvent = errors.New("no such event")
	ErrNoTask  = errors.New("no such task")
)

type Task struct {
	ID           string            `json:"id"`
	SubscribesTo []string          `json:"subscribesTo"`
	Image        string            `json:"image"`
	Env          map[string]string `json:"env"`
	CPU          float64           `json:"cpu"`
	Mem          float64           `json:"mem"`
}

type Event struct {
	ID      string            `json:"id"`
	Type    string            `json:"type"`
	Payload []byte            `json:"payload"`
	Meta    map[string]string `json:"meta"`
}

type TaskStorer interface {
	List() ([]*Task, error)
	Get(string) (*Task, error)
	Update(*Task) error
	Delete(string) error
}

type EventStorer interface {
	List() ([]*Event, error)
	Get(string) (*Event, error)
	Update(*Event) error
	Delete(string) error
}

type Storer interface {
	Tasks() TaskStorer
	Events() EventStorer
}
