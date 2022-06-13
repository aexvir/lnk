package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	gateway "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/aexvir/lnk/internal/logging"
	"github.com/aexvir/lnk/internal/storage"
	"github.com/aexvir/lnk/internal/svc"
	"github.com/aexvir/lnk/proto"
)

const port = 8000
const grpcaddr = ":9000"

func main() {
	log := logging.NewLogger("server")

	listener, err := net.Listen("tcp", grpcaddr)
	if err != nil {
		panic(err)
	}

	store, err := storage.NewMemoryStorage()
	if err != nil {
		panic(err)
	}

	grpcsrv := grpc.NewServer()
	linksvc := svc.NewLinksService(store)

	proto.RegisterLinksServer(grpcsrv, &linksvc)
	reflection.Register(grpcsrv)

	go func() {
		err := grpcsrv.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()

	apimux := gateway.NewServeMux()
	rpcopts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = proto.RegisterLinksHandlerFromEndpoint(context.Background(), apimux, grpcaddr, rpcopts)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	// todo: replace with different mux that allows more advanced routing
	mux.HandleFunc("/api/docs", svc.OpenapiDocsHandler)
	mux.HandleFunc("/api/schema.yaml", svc.OpenapiSchemaHandler)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api") {
			apimux.ServeHTTP(w, r)
			return
		}
		svc.LinkRedirectHandler(store)(w, r)
	})

	log.Write("startup", "listening on port %d", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	if err != nil {
		panic(err)
	}
}
