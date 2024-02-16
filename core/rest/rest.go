package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{w, r}
}

func (i *Context) Param(key string) string {
	return chi.URLParam(i.Request(), key)
}

func (i *Context) User() string {

	ctx, err := GetAuthContext(i)

	if err != nil {
		panic(ErrNotAuthorized)
	}

	return ctx.User
}

func (c *Context) ResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *Context) Request() *http.Request {
	return c.r
}

// response headers
func (i *Context) Headers() http.Header {
	return i.w.Header()
}

// get a header from request
func (i *Context) Header(key string) string {
	return i.r.Header.Get(key)
}

func (i *Context) Query(key string) string {
	return i.r.URL.Query().Get(key)
}

func (i *Context) Read(v any) (err error) {
	return json.NewDecoder(i.r.Body).Decode(v)
}

func (i *Context) Status(statusCode int) {
	i.w.WriteHeader(statusCode)
}

func (c *Context) Write(v any, statusCode ...int) {

	if len(statusCode) > 0 {
		c.w.WriteHeader(statusCode[0])
	}

	if str, isStr := v.(string); isStr {
		c.w.Write([]byte(str))
		return
	}

	if err, isErr := v.(error); isErr {
		c.w.Write([]byte(err.Error()))
		return
	}

	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)

	if err := enc.Encode(v); err != nil {
		http.Error(c.w, err.Error(), http.StatusInternalServerError)
		return
	}

	c.w.Header().Set("Content-Type", "application/json")

	c.w.Write(buf.Bytes())
}

type RestHandler func(c *Context)

type Service struct {
	http.Handler
	Name    string
	Version string
	mux     *chi.Mux
}

func NewService(name, version string) *Service {

	srv := &Service{
		Name:    name,
		Version: version,
		mux:     chi.NewMux(),
	}

	srv.Get("/version", func(c *Context) {
		c.ResponseWriter().Write([]byte(version))
	})

	return srv
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Service) Connect(pattern string, h RestHandler) {
	s.mux.Connect(pattern, wrap(h))
}

func (s *Service) Delete(pattern string, h RestHandler) {
	s.mux.Delete(pattern, wrap(h))
}

func (s *Service) Get(pattern string, h RestHandler) {
	s.mux.Get(pattern, wrap(h))
}

func (s *Service) Head(pattern string, h RestHandler) {
	s.mux.Head(pattern, wrap(h))
}

func (s *Service) Options(pattern string, h RestHandler) {
	s.mux.Options(pattern, wrap(h))
}

func (s *Service) Patch(pattern string, h RestHandler) {
	s.mux.Patch(pattern, wrap(h))
}

func (s *Service) Post(pattern string, h RestHandler) {
	s.mux.Post(pattern, wrap(h))
}

func (s *Service) Put(pattern string, h RestHandler) {
	s.mux.Put(pattern, wrap(h))
}

func (s *Service) Trace(pattern string, h RestHandler) {
	s.mux.Trace(pattern, wrap(h))
}

func wrap(h RestHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if msg := recover(); msg != nil {

				if err, isError := msg.(error); isError {
					if errors.Is(err, ErrNotAuthorized) {
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
				}

				slog.Error("rest:", msg)

				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		h(&Context{
			w: w,
			r: r,
		})
	}
}
