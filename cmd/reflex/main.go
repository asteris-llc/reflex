package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/asteris-llc/reflex/http"
	"github.com/asteris-llc/reflex/scheduler"
	"github.com/asteris-llc/reflex/state"
	"github.com/gogo/protobuf/proto"
	"github.com/hashicorp/consul/api"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
)

func startScheduler() (reflex *scheduler.ReflexScheduler, start func()) {
	exec := &mesos.ExecutorInfo{
		ExecutorId: util.NewExecutorID("default"),
		Name:       proto.String("BH executor (Go)"),
		Source:     proto.String("bh_test"),
		Command: &mesos.CommandInfo{
			Value: proto.String(""),
			Uris:  []*mesos.CommandInfo_URI{},
		},
	}

	fwinfo := &mesos.FrameworkInfo{
		User: proto.String(""),
		Name: proto.String("reflex"),
	}

	// skipping creds for now...
	// cred := (*mesos.Credential)(nil)
	// if *mesosAuthPrincipal != "" {
	// 	fwinfo.Principal = proto.String(*mesosAuthPrincipal)
	// 	secret, err := ioutil.ReadFile(*mesosAuthSecretFile)
	// 	if err != nil {
	// 		logrus.WithField("error", err).Fatal("failed reading secret file")
	// 	}
	// 	cred = &mesos.Credential{
	// 		Principal: proto.String(*mesosAuthPrincipal),
	// 		Secret:    secret,
	// 	}
	// }

	reflex = scheduler.NewScheduler(exec)

	config := sched.DriverConfig{
		Scheduler: reflex,
		Framework: fwinfo,
		Master:    "127.0.0.1:5050", // TODO: grab this from somewhere
		// Credential: cred,
	}

	start = func() {
		driver, err := sched.NewMesosSchedulerDriver(config)
		if err != nil {
			logrus.WithField("error", err).Fatal("unable to create a SchedulerDriver")
		}
		if stat, err := driver.Run(); err != nil {
			logrus.WithFields(logrus.Fields{
				"status": stat.String(),
				"error":  err,
			}).Info("framework stopped")
		}
	}

	return reflex, start
}

func newStore() *state.ConsulStore {
	client, err := api.NewClient(&api.Config{
		Address: "localhost:8500",
		Scheme:  "http",
	})
	if err != nil {
		panic(err)
	}

	return state.NewConsulStore("reflex", client)
}

func HandleNewEvents(store state.TaskStorer, events chan *state.Event, out chan *state.IOPair) {
	for {
		event := <-events
		logrus.WithField("event", event).Debug("new event")

		tasks, err := store.GetBySubscription(event.Type)
		if err != nil {
			panic(err) // TODO: make this not suck
		}

		for _, task := range tasks {
			out <- &state.IOPair{
				Event: event,
				Task:  task,
			}
		}
	}
}

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	store := newStore()

	api := http.NewAPI(store)
	go api.Serve("localhost:4000")

	reflex, start := startScheduler()

	go HandleNewEvents(store.Tasks(), api.Events, reflex.In)
	start()
}
