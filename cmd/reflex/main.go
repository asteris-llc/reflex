package main

import (
	"github.com/asteris-llc/reflex/http"
)

func main() {
	api := http.API{}
	api.Serve("localhost:4000")
}
