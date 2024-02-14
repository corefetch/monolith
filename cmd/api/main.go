package main

import (
	"net/http"
	"os"

	"learnt.io/core"
	"learnt.io/modules/accounts"
	"learnt.io/modules/payments"
	"learnt.io/modules/uploads"
	"learnt.io/modules/vcr"
	"learnt.io/modules/ws"
)

func main() {

	srv := core.Services()
	srv.Mount(accounts.Service())
	srv.Mount(uploads.Service())
	srv.Mount(payments.Service())
	srv.Mount(vcr.Service())
	srv.Mount(ws.Service())

	http.ListenAndServe(os.Getenv("LISTEN"), srv)
}
