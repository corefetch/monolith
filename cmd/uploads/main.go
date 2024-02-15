package main

import (
	"net/http"
	"os"

	"corefetch/modules/uploads"
)

func main() {
	http.ListenAndServe(os.Getenv("LISTEN"), uploads.Service())
}
