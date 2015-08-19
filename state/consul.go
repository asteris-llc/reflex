package state

import (
	"encoding/json"
	"github.com/hashicorp/consul/api"
)

type ConsulStore struct {
	root   string
	client *api.KV
}

func NewConsulStore(root string, client *api.Client) *ConsulStore {
	return &ConsulStore{root, client.KV()}
}

func (c *ConsulStore) inPath(path string) string {
	return c.root + "/" + path
}

func (c *ConsulStore) Tasks() TaskStorer {
	return &ConsulTaskStore{c}
}

func (c *ConsulStore) Events() EventStorer {
	return &ConsulEventStore{c}
}

// Tasks

type ConsulTaskStore struct {
	*ConsulStore
}

func (c *ConsulTaskStore) List() ([]*Task, error) {
	kvs, _, err := c.client.List(c.inPath("tasks"), new(api.QueryOptions))
	if err != nil {
		return nil, err
	}

	tasks := []*Task{}
	for _, kv := range kvs {
		task := new(Task)
		err := json.Unmarshal(kv.Value, &task)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (c *ConsulTaskStore) Get(id string) (*Task, error) {
	kv, _, err := c.client.Get(c.inPath("tasks/"+id), new(api.QueryOptions))
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, ErrNoTask
	}

	task := new(Task)
	err = json.Unmarshal(kv.Value, task)
	return task, err
}

func (c *ConsulTaskStore) Update(e *Task) error {
	// TODO: use check-and-set
	blob, err := json.Marshal(e)
	if err != nil {
		return err
	}
	kv := &api.KVPair{
		Key:   c.inPath("tasks/" + e.ID),
		Value: blob,
	}

	_, err = c.client.Put(kv, new(api.WriteOptions))
	return err
}

func (c *ConsulTaskStore) Delete(id string) error {
	_, err := c.client.Delete(c.inPath("tasks/"+id), new(api.WriteOptions))
	return err
}

// Events

type ConsulEventStore struct {
	*ConsulStore
}

func (c *ConsulEventStore) List() ([]*Event, error) {
	kvs, _, err := c.client.List(c.inPath("events"), new(api.QueryOptions))
	if err != nil {
		return nil, err
	}

	events := []*Event{}
	for _, kv := range kvs {
		event := new(Event)
		err := json.Unmarshal(kv.Value, &event)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (c *ConsulEventStore) Get(id string) (*Event, error) {
	kv, _, err := c.client.Get(c.inPath("events/"+id), new(api.QueryOptions))
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, ErrNoEvent
	}

	event := new(Event)
	err = json.Unmarshal(kv.Value, event)
	return event, err
}

func (c *ConsulEventStore) Update(e *Event) error {
	// TODO: use check-and-set
	blob, err := json.Marshal(e)
	if err != nil {
		return err
	}
	kv := &api.KVPair{
		Key:   c.inPath("events/" + e.ID),
		Value: blob,
	}

	_, err = c.client.Put(kv, new(api.WriteOptions))
	return err
}

func (c *ConsulEventStore) Delete(id string) error {
	_, err := c.client.Delete(c.inPath("events/"+id), new(api.WriteOptions))
	return err
}
