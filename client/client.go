package main

import (
	"context"
	"flag"
	"log"
	"time"

	pb "github.com/net-reply-future-networks/k8s-golang-append-only-store/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var addr = flag.String("addr", ":3000", "the address to connect to")

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDatastoreClient(conn)

	// set(c, "alex", "brown")
	// set(c, "micheal", "green")
	// set(c, "alex", "brown")
	// set(c, "micheal", "green")
	// set(c, "alex", "brown")
	// set(c, "micheal", "green")
	get(c, "alex")
	// del(c, "alex")
}

// SET tombstone, key, value
func set(c pb.DatastoreClient, key string, value string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Set(ctx, &pb.SetRequest{Key: key, Value: value})
	if err != nil {
		log.Printf("set request denied. error description: %v", err)
	}
	log.Printf("Return Message: %s", r.Message)
}

// GET key, value
func get(c pb.DatastoreClient, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Get(ctx, &pb.GetRequest{Key: key})
	if err != nil {
		log.Fatalf("could not make get request for key: %v", err)
	}
	log.Printf("Return Data: %v %v", r.Key, r.Value)
}

// DEL key, value
func del(c pb.DatastoreClient, key string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Del(ctx, &pb.DelRequest{Key: key})
	if err != nil {
		log.Fatalf("could not make delete request for key: %v", err)
	}
	log.Printf("Return Data: %v", r.Message)
}
