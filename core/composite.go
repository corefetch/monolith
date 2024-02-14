package core

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"learnt.io/core/rest"
)

type services struct {
	mux chi.Router
}

func Services() *services {
	return &services{
		mux: chi.NewMux(),
	}
}

func (s *services) Mount(srv *rest.Service) {
	s.mux.Mount("/"+srv.Name, srv)
}

func (s *services) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
