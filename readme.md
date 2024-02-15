# CoreFetch

A ready to use monolith flexible to decouple in microservices for bootstraping your backend application in GoLang.

```
package main

import (
	"net/http"
	"os"

	"corefetch/core/rest"
	"corefetch/modules/accounts"
	"corefetch/modules/payments"
	"corefetch/modules/uploads"
	"corefetch/modules/vcr"
	"corefetch/modules/ws"
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
```