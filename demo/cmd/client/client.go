package main

import (
	"context"
	api "do/api/pinger"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial grpc: %v", err)
	}
	defer conn.Close()

	testPing(conn)
	// testLostConn(conn)
	// testClientStream(conn)
	// testServerStream(conn)
	// testBidiStream(conn)
}

func testPing(conn *grpc.ClientConn) {
	client := api.NewPingerClient(conn)
	r, err := client.Ping(context.Background(), &api.PingRequest{Message: "ping" + time.Now().String()})
	if err != nil {
		log.Fatalf("could not ping: %v", err)
	}
	fmt.Printf("ping reply: %v\n", r)
}

func testLostConn(conn *grpc.ClientConn) {
	client := api.NewPingerClient(conn)
	for {
		r, err := client.Ping(context.Background(), &api.PingRequest{Message: "ping" + time.Now().String()})
		if err != nil {
			fmt.Printf("could not ping: %v", err)
		}
		fmt.Printf("ping reply: %v\n", r)
		time.Sleep(time.Second)
	}
}

func testBidiStream(conn *grpc.ClientConn) {
	client := api.NewPingerClient(conn)
	ctx := context.Background()
	stream, err := client.BidiStream(ctx)
	if err != nil {
		log.Fatalf("could not get server stream: %v", err)
	}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("end of stream")
			break
		}
		if err != nil {
			log.Fatalf("could not handle stream data: %v", err)
		}
		log.Println(data)
		if err := stream.Send(&api.StreamRequest{Message: fmt.Sprintf("client ping %d", time.Now().Unix())}); err != nil {
			log.Fatalf("could not send client stream: %v", err)
		}
	}
}

func testClientStream(conn *grpc.ClientConn) {
	client := api.NewPingerClient(conn)
	ctx := context.Background()
	stream, err := client.ClientStream(ctx)
	if err != nil {
		log.Fatalf("could not get server stream: %v", err)
	}
	for {
		if err := stream.Send(&api.StreamRequest{Message: fmt.Sprintf("client ping %d", time.Now().Unix())}); err != nil {
			log.Printf("could not send client stream: %v", err)
		}
		time.Sleep(time.Second)
	}
}

func testServerStream(conn *grpc.ClientConn) {
	client := api.NewPingerClient(conn)
	ctx := context.Background()
	stream, err := client.ServerStream(ctx, &api.StreamRequest{Message: "initial client request"})
	if err != nil {
		log.Fatalf("could not get server stream: %v", err)
	}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("end of stream")
			break
		}
		log.Println(data)
		if err != nil {
			fmt.Printf("could not handle stream data: %v", err)
		}
	}
}
