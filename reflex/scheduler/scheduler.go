package scheduler

import (
	"github.com/Sirupsen/logrus"
	// "github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	// util "github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
)

type ReflexScheduler struct {
	executor *mesos.ExecutorInfo
}

func NewScheduler(exec *mesos.ExecutorInfo) *ReflexScheduler {
	sched := &ReflexScheduler{
		executor: exec,
	}

	return sched
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
				cpus += res.GetScalar().GetValue()
			}
		}

		logrus.WithFields(logrus.Fields{
			"cpus": cpus,
			"mem":  mem,
		}).Debug("got offer")

		tasks := []*mesos.TaskInfo{}

		// TODO: get a list of things that should be running from logic, and try to
		// schedule some of them. Let the status messages reflect the actual run
		// state back to logic.

		// TODO: move this to logic
		// for _, pair := range sched.waitingPairs {
		// 	if pair.InProgress {
		// 		continue
		// 	}

		// 	task := pair.Task
		// 	event := pair.Event

		// 	if cpus >= task.CPU && mem >= task.Mem {

		// 		info := &mesos.TaskInfo{
		// 			TaskId: &mesos.TaskID{
		// 				Value: proto.String("reflex-" + event.ID),
		// 			},
		// 			Name:    proto.String("EXEC_reflex-" + event.ID),
		// 			SlaveId: offer.SlaveId,
		// 			Resources: []*mesos.Resource{
		// 				util.NewScalarResource("cpus", task.CPU),
		// 				util.NewScalarResource("mem", task.Mem),
		// 			},
		// 			Executor: &mesos.ExecutorInfo{
		// 				ExecutorId: &mesos.ExecutorID{Value: proto.String("reflex-executor")},
		// 				Command: &mesos.CommandInfo{
		// 					Value: proto.String("cat"),
		// 				},
		// 				Name: proto.String("reflex"),
		// 			},
		// 			Data: event.Payload,
		// 		}

		// 		tasks = append(tasks, info)

		// 		cpus -= task.CPU
		// 		mem -= task.CPU
		// 	}

		// 	if cpus <= 0 || mem <= 0 {
		// 		break
		// 	}
		// }

		filters := new(mesos.Filters)

		if len(tasks) == 0 {
			driver.DeclineOffer(offer.GetId(), filters)
		} else {
			driver.LaunchTasks(
				[]*mesos.OfferID{offer.GetId()},
				tasks,
				filters,
			)
		}
	}
}

func (sched *ReflexScheduler) StatusUpdate(driver sched.SchedulerDriver, status *mesos.TaskStatus) {
	// TODO: extract the ID and send these task updates to logic, along with a
	// "nice" version of the task status.

	logrus.WithFields(logrus.Fields{
		"status": status, // TODO: parse these fields out so it's not such a mess
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
