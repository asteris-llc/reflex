package main

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/asteris-llc/reflex/reflex/logic"
	exec "github.com/mesos/mesos-go/executor"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/samalba/dockerclient"
	"os"
)

type ReflexExecutor struct {
	data *logic.IOPair
}

func newReflexExecutor() *ReflexExecutor {
	return &ReflexExecutor{}
}

func (exec *ReflexExecutor) Registered(driver exec.ExecutorDriver, execInfo *mesos.ExecutorInfo, fwinfo *mesos.FrameworkInfo, slaveInfo *mesos.SlaveInfo) {
	data := new(logic.IOPair)
	err := json.Unmarshal(execInfo.Data, data)
	if err != nil {
		logrus.WithField("error", err).Fatal("could not parse task info")
	}

	exec.data = data

	logrus.WithField("host", slaveInfo.GetHostname()).Info("registered executor")
}

func (exec *ReflexExecutor) Reregistered(driver exec.ExecutorDriver, slaveInfo *mesos.SlaveInfo) {
	logrus.WithField("host", slaveInfo.GetHostname()).Info("re-registered executor")
}

func (exec *ReflexExecutor) Disconnected(exec.ExecutorDriver) {
	logrus.Warning("executor disconnected")
}

func (exec *ReflexExecutor) LaunchTask(driver exec.ExecutorDriver, taskInfo *mesos.TaskInfo) {
	logrus.WithFields(logrus.Fields{
		"task":    taskInfo.GetName(),
		"command": taskInfo.Command.GetValue(), // TODO: probably the image instead
	}).Info("launching task")

	runStatus := &mesos.TaskStatus{
		TaskId: taskInfo.GetTaskId(),
		State:  mesos.TaskState_TASK_RUNNING.Enum(),
	}
	_, err := driver.SendStatusUpdate(runStatus)
	if err != nil {
		logrus.WithField("error", err).Error("got error sending status update")
	}

	// start the docker container
	image := exec.data.Task.Image
	// payload := exec.data.Event.Payload

	client, err := dockerclient.NewDockerClient(os.Getenv("DOCKER_HOST"), nil)
	if err != nil {
		logrus.WithField("error", err).Error("could not start docker client")
		// TODO: send a failed message and exit
	}

	containerId, err := client.CreateContainer(
		&dockerclient.ContainerConfig{
			AttachStdin: true,
			Image:       image,
		},
		taskInfo.GetTaskId().String(),
	)
	if err != nil {
		logrus.WithField("error", err).Error("could not create the docker container")
	}

	err = client.StartContainer(containerId, new(dockerclient.HostConfig))
	if err != nil {
		logrus.WithField("error", err).Error("could not start the docker container")
	}

	// TODO: move the below to a separate method for use when things fail during startup
	// finish task
	fmt.Println("Finishing task", taskInfo.GetName())

	finStatus := &mesos.TaskStatus{
		TaskId: taskInfo.GetTaskId(),
		State:  mesos.TaskState_TASK_FINISHED.Enum(),
	}
	_, err = driver.SendStatusUpdate(finStatus)
	if err != nil {
		logrus.WithField("error", err).Error("got error sending status update")
	}

	logrus.WithField("task", taskInfo.GetName()).Info("task finished")
}

func (exec *ReflexExecutor) KillTask(exec.ExecutorDriver, *mesos.TaskID) {
	logrus.Info("kill task")
}

func (exec *ReflexExecutor) FrameworkMessage(driver exec.ExecutorDriver, msg string) {
	logrus.WithField("message", msg).Info("got framework message")
}

func (exec *ReflexExecutor) Shutdown(exec.ExecutorDriver) {
	logrus.Info("shutting down the executor")
}

func (exec *ReflexExecutor) Error(driver exec.ExecutorDriver, err string) {
	logrus.WithField("error", err).Error("got error message")
}
