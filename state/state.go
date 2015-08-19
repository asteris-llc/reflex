package state

import (
	"errors"
)

var (
	ErrNoEvent = errors.New("no such event")
	ErrNoTask  = errors.New("no such task")
)

type Task struct {
	ID string
}

type Event struct {
	ID string
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
