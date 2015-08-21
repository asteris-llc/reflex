package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/asteris-llc/reflex/reflex"
)

func main() {
	logrus.SetLevel(logrus.DebugLevel)

	r, err := reflex.New(&reflex.Options{
		Address: ":4000",
	})
	if err != nil {
		panic(err)
	}

	r.Start()
}
