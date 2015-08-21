package main

import (
	"github.com/Sirupsen/logrus"
	exec "github.com/mesos/mesos-go/executor"
)

// TODO: move this functionality underneath `cmd`
func main() {
	logrus.Info("starting reflex executor")

	driver, err := exec.NewMesosExecutorDriver(
		exec.DriverConfig{
			Executor: newReflexExecutor(),
		},
	)
	if err != nil {
		logrus.WithField("error", err).Fatal("could not create executor")
	}

	_, err = driver.Start()
	if err != nil {
		logrus.WithField("error", err).Fatal("could not start executor")
	}

	logrus.Info("started")
	driver.Join()
}
