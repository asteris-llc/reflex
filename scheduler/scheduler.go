package scheduler

import (
	"github.com/Sirupsen/logrus"
	"github.com/asteris-llc/reflex/state"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
	"golang.org/x/net/context"
	"sync"
)

type ReflexScheduler struct {
	executor *mesos.ExecutorInfo

	In  chan *state.Task
	Out chan *state.Task

	// task pool
	taskLock     sync.Mutex
	waitingTasks []*state.Task

	// meta-state
	context context.Context
	cancel  func()
}

func NewScheduler(exec *mesos.ExecutorInfo) *ReflexScheduler {
	context, cancel := context.WithCancel(context.Background())

	sched := &ReflexScheduler{
		executor: exec,
		In:       make(chan *state.Task),
		Out:      make(chan *state.Task, 100), // TODO: make this buffer configurable

		taskLock:     sync.Mutex{},
		waitingTasks: []*state.Task{},

		context: context,
		cancel:  cancel,
	}

	go sched.slurpTasks()

	return sched
}

func (sched *ReflexScheduler) slurpTasks() {
	for {
		select {
		case task := <-sched.In:
			sched.waitingTasks = append(sched.waitingTasks, task)
		case <-sched.context.Done():
			logrus.Info("stopping task slurper")
			return
		}
	}
}

func (sched *ReflexScheduler) Stop() {
	// TODO: this should stop the scheduler too
	sched.cancel()
}

func (sched *ReflexScheduler) Registered(driver sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	logrus.WithFields(logrus.Fields{
		"masterInfo": masterInfo,
	}).Info("registered with master")
}

func (sched *ReflexScheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	logrus.WithFields(logrus.Fields{
		"masterInfo": masterInfo,
	}).Info("re-registered with master")
}

func (sched *ReflexScheduler) Disconnected(sched.SchedulerDriver) {
	logrus.Info("disconnected")
}

func (sched *ReflexScheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {
	for _, offer := range offers {
		// CPUs
		cpuResources := util.FilterResources(
			offer.Resources,
			func(res *mesos.Resource) bool {
				return res.GetName() == "cpus"
			},
		)
		cpus := 0.0
		for _, res := range cpuResources {
			cpus += res.GetScalar().GetValue()
		}

		// Mem
		memResources := util.FilterResources(
			offer.Resources,
			func(res *mesos.Resource) bool {
				return res.GetName() == "mem"
			},
		)
		mem := 0.0
		for _, res := range memResources {
			mem += res.GetScalar().GetValue()
		}

		logrus.WithFields(logrus.Fields{
			"cpus": cpus,
			"mem":  mem,
		}).Info("got offer")
	}
}

func (sched *ReflexScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	logrus.WithFields(logrus.Fields{
		"task":  status.TaskId.GetValue(),
		"state": status.State.Enum().String(),
	}).Info("got status update")
}

func (sched *ReflexScheduler) OfferRescinded(sched.SchedulerDriver, *mesos.OfferID) {}
func (sched *ReflexScheduler) FrameworkMessage(sched.SchedulerDriver, *mesos.ExecutorID, *mesos.SlaveID, string) {
}
func (sched *ReflexScheduler) SlaveLost(sched.SchedulerDriver, *mesos.SlaveID) {}
func (sched *ReflexScheduler) ExecutorLost(sched.SchedulerDriver, *mesos.ExecutorID, *mesos.SlaveID, int) {
}

func (sched *ReflexScheduler) Error(driver sched.SchedulerDriver, err string) {
	logrus.WithField("error", err).Error("got error")
}
