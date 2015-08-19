package main

import (
	"github.com/asteris-llc/reflex/http"
	"github.com/asteris-llc/reflex/state"
	"github.com/hashicorp/consul/api"
)

func main() {
	client, err := api.NewClient(&api.Config{
		Address: "localhost:8500",
		Scheme:  "http",
	})
	if err != nil {
		panic(err)
	}

	store := state.NewConsulStore("reflex", client)

	api := http.NewAPI(store)
	api.Serve("localhost:4000")
}
