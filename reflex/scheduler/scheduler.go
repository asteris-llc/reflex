package scheduler

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"github.com/asteris-llc/reflex/reflex/logic"
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
)

type ReflexScheduler struct {
	executor *mesos.ExecutorInfo
	filters  *mesos.Filters
	logic    *logic.Logic
}

func NewScheduler(exec *mesos.ExecutorInfo, logic *logic.Logic) (*ReflexScheduler, error) {
	sched := &ReflexScheduler{
		executor: exec,
		filters:  new(mesos.Filters), // TODO: make timeout tunable
		logic:    logic,
	}

	return sched, nil
}

func (r *ReflexScheduler) Start(fwinfo *mesos.FrameworkInfo, master string) {
	config := sched.DriverConfig{
		Scheduler: r,
		Framework: fwinfo,
		Master:    master,
		// TODO: Credential: cred,
	}

	driver, err := sched.NewMesosSchedulerDriver(config)
	if err != nil {
		logrus.WithField("error", err).Fatal("could not create MesosSchedulerDriver")
		return
	}

	if stat, err := driver.Run(); err != nil {
		logrus.WithFields(logrus.Fields{
			"status": stat.String(),
			"error":  err,
		}).Fatal("framework stopped")
	}
}

// TODO: where does reconciliation go here?

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
		cpus := 0.0
		mem := 0.0

		for _, res := range offer.Resources {
			switch res.GetName() {
			case "cpus":
				cpus += res.GetScalar().GetValue()
			case "mem":
				mem += res.GetScalar().GetValue()
			}
		}

		logrus.WithFields(logrus.Fields{
			"cpus": cpus,
			"mem":  mem,
		}).Debug("got offer")

		tasks := []*mesos.TaskInfo{}

		for _, pair := range sched.logic.ToSchedule() {
			logrus.WithField("id", pair.ID).Debug("scheduling task")
			task := pair.Task

			// stop early if the offer isn't big enough
			if cpus <= task.CPU || mem <= task.Mem {
				logrus.WithFields(logrus.Fields{
					"pair":       pair,
					"cpusReq":    task.CPU,
					"cpusActual": cpus,
					"memReq":     task.Mem,
					"memActual":  mem,
				}).Debug("remaining offer not big enough")
				continue
			}

			payload, err := json.Marshal(pair)
			if err != nil {
				panic(err) // TODO: handle this more gracefully
			}

			info := &mesos.TaskInfo{
				TaskId: &mesos.TaskID{
					Value: proto.String(pair.ID),
				},
				Name:    proto.String("EXEC_reflex-" + pair.ID),
				SlaveId: offer.SlaveId,
				Resources: []*mesos.Resource{
					util.NewScalarResource("cpus", task.CPU),
					util.NewScalarResource("mem", task.Mem),
				},
				Executor: &mesos.ExecutorInfo{
					ExecutorId: &mesos.ExecutorID{Value: proto.String("reflex-exeutor")},
					Name:       proto.String("reflex-executor"),
					Command: &mesos.CommandInfo{
						Value: proto.String("asdfasdfasdf"),
						Uris:  []*mesos.CommandInfo_URI{},
					},
					Data: payload,
				},
			}

			tasks = append(tasks, info)
			sched.logic.TaskStarted(pair.ID)

			cpus -= task.CPU
			mem -= task.CPU

			if cpus <= 0 || mem <= 0 {
				break
			}
		}

		if len(tasks) == 0 {
			driver.DeclineOffer(offer.GetId(), sched.filters) // TODO: handle error
		} else {
			_, err := driver.LaunchTasks(
				[]*mesos.OfferID{offer.GetId()},
				tasks,
				sched.filters,
			)
			if err != nil {
				panic(err) // TODO: handle this more gracefully
			}
		}
	}
}

func (sched *ReflexScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	logrus.WithFields(logrus.Fields{
		"status": status, // TODO: parse these fields out so it's not such a mess
	}).Info("got status update")

	switch *status.State {
	case mesos.TaskState_TASK_STAGING, mesos.TaskState_TASK_STARTING, mesos.TaskState_TASK_RUNNING:
		sched.logic.TaskStarted(*status.TaskId.Value)
	case mesos.TaskState_TASK_FAILED, mesos.TaskState_TASK_ERROR, mesos.TaskState_TASK_KILLED, mesos.TaskState_TASK_LOST: // IE: only failure terminal states
		sched.logic.TaskFinished(*status.TaskId.Value, false)
	case mesos.TaskState_TASK_FINISHED: // IE: only successful terminal states
		sched.logic.TaskFinished(*status.TaskId.Value, true)
	}
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
