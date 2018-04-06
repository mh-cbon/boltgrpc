# Bolt gRPC

gRPC interface for boltdb

## Installation

    go get github.com/sirait/boltgrpc/...

## Run the server

    boltgrpc --db bolt.db --grpcport 9090 --httpport 9091

## Client example

```
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/sirait/boltgrpc"

	"google.golang.org/grpc"
)

func main() {
	// connect
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":9090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %s", err)
	}
	defer conn.Close()

	// create grpc client
	c := boltgrpc.NewBoltClient(conn)

	// update
	_, err = c.Update(context.Background(), &boltgrpc.UpdateRequest{Buckets: []string{"notes"}, Key: []byte("1"), Val: []byte("I'm trying boltgrpc")})
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	// view
	response, err := c.View(context.Background(), &boltgrpc.ViewRequest{Buckets: []string{"notes"}, Key: []byte("1")})
	if err != nil {
		log.Fatalf("Error: %s", err)
	}
	fmt.Println(string(response.Val))

}

```

## Backup

    curl http://localhost:9091/backup > backup.db

