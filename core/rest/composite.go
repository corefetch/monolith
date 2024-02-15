package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type services struct {
	mux chi.Router
}

func Services() *services {
	return &services{
		mux: chi.NewMux(),
	}
}

func (s *services) Use(srv *Service) {
	s.mux.Mount("/"+srv.Name, srv)
}

func (s *services) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
