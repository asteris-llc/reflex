package reflex

import (
	"github.com/asteris-llc/reflex/reflex/http"
	"github.com/asteris-llc/reflex/reflex/logic"
	"github.com/asteris-llc/reflex/reflex/scheduler"
	"github.com/kardianos/osext"
	"golang.org/x/net/context"
)

type Reflex struct {
	opts *Options

	context context.Context
	cancel  func()
}

type Options struct {
	Address string
}

func New(opts *Options) (*Reflex, error) {
	context, cancel := context.WithCancel(context.Background())
	return &Reflex{
		opts:    opts,
		context: context,
		cancel:  cancel,
	}, nil
}

func (r *Reflex) Start() error {
	events := make(chan *logic.Event)

	// HTTP
	api, err := http.NewAPI(events)
	if err != nil {
		return err
	}

	local, err := osext.ExecutableFolder()
	if err != nil {
		return err
	}
	artifacts, err := http.NewArtifacts(local)
	if err != nil {
		return err
	}

	http := http.HTTP{
		Components: []http.Registerer{api, artifacts},
	}
	go http.ServeHTTP(r.opts.Address) // TODO: figure out how to stop this

	// Logic and Scheduler
	logic, err := logic.NewLogic(r.context, events)
	if err != nil {
		return err
	}
	logic.Start()

	// Scheduler
	exec, fwinfo := mesosExecutorMeta()
	sched, err := scheduler.NewScheduler(exec, logic)
	go sched.Start(fwinfo, "localhost:5050")

	<-r.context.Done()

	return nil
}

func (r *Reflex) Stop() {
	r.cancel()
}
