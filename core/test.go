package core

import (
	"net/http"
	"os"

	"corefetch/core/rest"
)

func TestSetup() {
	os.Setenv("LISTEN", ":8888")
	os.Setenv("DB", "mongodb://learnt:learnt@localhost:27017/learnt_test")
	os.Setenv("UPLOADS_DIR", "/Users/cosmin/Work/Learn/backend/_uploads")
}

func TestService(s *rest.Service) {

	TestSetup()

	go func() {
		if err := http.ListenAndServe(":8888", s); err != nil {
			panic(err)
		}
	}()
}
