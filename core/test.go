package core

import (
	"net/http"
	"os"

	"learnt.io/core/rest"
)

func TestService(s *rest.Service) {
	os.Setenv("LISTEN", ":8888")
	os.Setenv("DB", "mongodb://learnt:learnt@localhost:27017/learnt_test")
	os.Setenv("UPLOADS_DIR", "/Users/cosmin/Work/Learn/backend/_uploads")
	go func() {
		if err := http.ListenAndServe(":8888", s); err != nil {
			panic(err)
		}
	}()
}
