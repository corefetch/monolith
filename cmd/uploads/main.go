package main

import (
	"net/http"
	"os"

	"learnt.io/modules/uploads"
)

func main() {
	http.ListenAndServe(os.Getenv("LISTEN"), uploads.Service())
}
