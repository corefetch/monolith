package main

import (
	"net/http"
	"os"

	"learnt.io/core/rest"
	"learnt.io/modules/accounts"
	"learnt.io/modules/payments"
	"learnt.io/modules/uploads"
	"learnt.io/modules/vcr"
	"learnt.io/modules/ws"
)

func main() {

	srv := rest.Services()
	srv.Use(accounts.Service())
	srv.Use(uploads.Service())
	srv.Use(payments.Service())
	srv.Use(vcr.Service())
	srv.Use(ws.Service())

	http.ListenAndServe(os.Getenv("LISTEN"), srv)
}
