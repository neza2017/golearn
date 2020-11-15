// generate code
// protoc -I proto  --go_out=plugins=grpc:proto proto/helloworld.proto
package main

import (
	"context"
	"fmt"
	pb "golearn/grpcExp/proto"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const (
	port    = ":50051"
	address = "localhost:50051"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func startServer() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	return s.Serve(lis)
}

func hello(msg string) error {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return err
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	ctx, cancle := context.WithTimeout(context.Background(), time.Second)
	defer cancle()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: msg})
	if err != nil {
		return err
	}
	log.Printf("Greeting: %s", r.GetMessage())
	return nil
}

func main() {
	go func() {
		if err := startServer(); err != nil {
			log.Fatal(err)
		}
	}()
	for i := 0; i < 10; i++ {
		if err := hello(fmt.Sprintf("-%d-", i)); err != nil {
			log.Fatal(err)
		}
	}

}
