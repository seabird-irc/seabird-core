//go:generate protoc -I ../../pb --go_out=plugins=grpc:../../pb/ ../../pb/seabird.proto

package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/belak/seabird-core/pb"
)

const (
	address = "localhost:50052"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSeabirdClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.Connect(ctx, &pb.ConnectRequest{})
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	log.Printf("Greeting: %s", r.GetClientId())

	stream, err := c.EventStream(ctx, &pb.EventStreamRequest{})
	if err != nil {
		log.Fatalf("could not get event stream: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("could not get event: %v", err)
		}

		log.Printf("Msg: %v", msg)
	}

	resp, err := c.SendMessage(ctx, &pb.SendMessageRequest{Channel: "#minecraft", Message: "hello world"})
	if err != nil {
		log.Fatalf("could not send message: %v", err)
	}

	log.Printf("Resp: %v\n", resp)
}
