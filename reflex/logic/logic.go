package logic

import (
	"github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"sync"
)

// Logic holds the "busines logic" for event handling and events around the
// cluster.
type Logic struct {
	events chan *Event

	pairs     map[string]*IOPair
	pairsLock *sync.RWMutex

	ctx context.Context
}

func NewLogic(ctx context.Context, events chan *Event) (*Logic, error) {
	logic := &Logic{
		events:    events,
		pairs:     make(map[string]*IOPair),
		pairsLock: new(sync.RWMutex),
		ctx:       ctx,
	}

	return logic, nil
}

func (l *Logic) Start() {
	go l.MakePairs(l.ctx)
}

func (l *Logic) MakePairs(ctx context.Context) {
	for {
		select {
		case event := <-l.events:
			logrus.WithField("event", event).Debug("received new event")
			l.addEvent(event)
		case <-ctx.Done():
			return
		}
	}
}

func (l *Logic) addEvent(e *Event) {
	l.pairsLock.Lock()
	defer l.pairsLock.Unlock()

	// TODO: add a new pair for each Task
	task := &Task{
		ID:    "test",
		Image: "brianhicks/reflex-stderr",
		Env: map[string]string{
			"TEST": "1",
		},
		CPU: 0.5,
		Mem: 256,
	}

	pair := &IOPair{
		ID:        e.ID, // TODO: should this be a new ID? I don't think so.
		Task:      task,
		Event:     e,
		Scheduled: false,
	}
	l.pairs[pair.ID] = pair
}

// ToSchedule is called by the scheduler when it has task offers to schedule
// work onto. It should return the most current list of IOPairs to schedule.
func (l *Logic) ToSchedule() []*IOPair {
	l.pairsLock.RLock()
	defer l.pairsLock.RUnlock()

	pairs := []*IOPair{}
	for _, pair := range l.pairs {
		if !pair.Scheduled {
			pairs = append(pairs, pair)
		}
	}

	return pairs
}

// TaskStarted is called by the framework when the task is started. It should be
// idempotent, because it might be called more than once.
func (l *Logic) TaskStarted(id string) {
	l.pairsLock.Lock()
	defer l.pairsLock.Unlock()

	logrus.WithField("id", id).Debug("task started")
	pair, ok := l.pairs[id]
	if !ok { // race condition. Must have just failed instead of starting properly
		return
	}

	pair.Scheduled = true
}

// TaskFinished is called by the framework when the task is finished (either
// lost or successful). It doesn't have to be idempotent.
func (l *Logic) TaskFinished(id string, success bool) {
	l.pairsLock.Lock()
	defer l.pairsLock.Unlock()

	base := logrus.WithFields(logrus.Fields{
		"success": success,
		"id":      id,
	})
	if success {
		base.Debug("task finished")
	} else {
		base.Error("task exited unsuccessfully")
	}

	delete(l.pairs, id)
}
