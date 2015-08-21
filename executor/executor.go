package main

import (
	"fmt"
	exec "github.com/mesos/mesos-go/executor"
	mesos "github.com/mesos/mesos-go/mesosproto"
)

type ReflexExecutor struct {
	data []byte
}

func newReflexExecutor() *ReflexExecutor {
	return &ReflexExecutor{}
}

func (exec *ReflexExecutor) Registered(driver exec.ExecutorDriver, execInfo *mesos.ExecutorInfo, fwinfo *mesos.FrameworkInfo, slaveInfo *mesos.SlaveInfo) {
	exec.data = execInfo.Data
	fmt.Println("Registered Executor on slave ", slaveInfo.GetHostname())
}

func (exec *ReflexExecutor) Reregistered(driver exec.ExecutorDriver, slaveInfo *mesos.SlaveInfo) {
	fmt.Println("Re-registered Executor on slave ", slaveInfo.GetHostname())
}

func (exec *ReflexExecutor) Disconnected(exec.ExecutorDriver) {
	fmt.Println("Executor disconnected.")
}

func (exec *ReflexExecutor) LaunchTask(driver exec.ExecutorDriver, taskInfo *mesos.TaskInfo) {
	fmt.Println("Launching task", taskInfo.GetName(), "with command", taskInfo.Command.GetValue())

	runStatus := &mesos.TaskStatus{
		TaskId: taskInfo.GetTaskId(),
		State:  mesos.TaskState_TASK_RUNNING.Enum(),
	}
	_, err := driver.SendStatusUpdate(runStatus)
	if err != nil {
		fmt.Println("Got error", err)
	}

	//
	// this is where one would perform the requested task
	//

	// finish task
	fmt.Println("Finishing task", taskInfo.GetName())

	finStatus := &mesos.TaskStatus{
		TaskId: taskInfo.GetTaskId(),
		State:  mesos.TaskState_TASK_FINISHED.Enum(),
	}
	_, err = driver.SendStatusUpdate(finStatus)
	if err != nil {
		fmt.Println("Got error", err)
	}
	fmt.Println("Task finished", taskInfo.GetName())
}

func (exec *ReflexExecutor) KillTask(exec.ExecutorDriver, *mesos.TaskID) {
	fmt.Println("Kill task")
}

func (exec *ReflexExecutor) FrameworkMessage(driver exec.ExecutorDriver, msg string) {
	fmt.Println("Got framework message: ", msg)
}

func (exec *ReflexExecutor) Shutdown(exec.ExecutorDriver) {
	fmt.Println("Shutting down the executor")
}

func (exec *ReflexExecutor) Error(driver exec.ExecutorDriver, err string) {
	fmt.Println("Got error message:", err)
}
