package main

import (
	"fmt"
	"net/http"

	"github.com/aexvir/lnk/internal/logging"
	"github.com/aexvir/lnk/internal/storage"
	"github.com/aexvir/lnk/internal/svc"
)

const port = 8000

func main() {
	log := logging.NewLogger("server")

	store, err := storage.NewMemoryStorage()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// todo: replace with different mux that allows more advanced routing
	mux.HandleFunc("/api/docs", svc.OpenapiDocsHandler)
	mux.HandleFunc("/api/schema.json", svc.OpenapiSchemaHandler)
	mux.HandleFunc("/api/links/", svc.LinkMgmtHandler(store))
	mux.HandleFunc("/api/links", svc.LinkMgmtHandler(store))

	mux.HandleFunc("/", svc.LinkRedirectHandler(store))

	log.Write("startup", "listening on port %d", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		panic(err)
	}
}
