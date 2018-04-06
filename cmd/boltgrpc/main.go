package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/sirait/boltgrpc"

	"google.golang.org/grpc"
)

func main() {

	// get flags
	path := flag.String("db", "bolt.db", "db file path")
	grpcPort := flag.String("grpcport", "9090", "grpc port")
	httpPort := flag.String("httpport", "9091", "http port")
	flag.Parse()

	// create handler and server
	boltHandler := boltgrpc.Handler{*path}
	boltServer := grpc.NewServer()

	// register bolt server
	boltgrpc.RegisterBoltServer(boltServer, &boltHandler)

	// run http server
	go func() {
		http.HandleFunc("/backup", func(w http.ResponseWriter, r *http.Request) {
			boltHandler.Backup(w, r)
		})

		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *httpPort), nil))
	}()

	// run bolt grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err := boltServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
