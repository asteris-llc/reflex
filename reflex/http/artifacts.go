package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Artifacts struct {
	ArtifactPath string
}

func NewArtifacts(path string) (*Artifacts, error) {
	return &Artifacts{path}, nil
}

func (a *Artifacts) Register(router *mux.Router) {
	webPath := "/artifacts"

	// TODO: maybe use http.ServeFile for security's sake. This is good enough for
	// now, though.
	router.PathPrefix(webPath).Handler(
		http.StripPrefix(webPath, http.FileServer(http.Dir(a.ArtifactPath))),
	)
}
