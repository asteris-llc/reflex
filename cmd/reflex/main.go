package main

import (
	"github.com/asteris-llc/reflex/reflex"
)

func main() {
	r, err := reflex.New(&reflex.Options{
		Address: ":4000",
	})
	if err != nil {
		panic(err)
	}

	r.Start()
}
