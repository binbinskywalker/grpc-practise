package main

import (
	"context"
	"fmt"
	"net"

	pb "./../../proto"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	Address = "localhost:8888"
)

type helloService struct{}

var HelloService = helloService{}

func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponce, error) {
	resp := new(pb.HelloResponce)
	resp.message = fmt.Sprintf("Hello %s,", in.Name)
	return resp, nil
}
func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("fail to listen :%v", err)
	}
	s := grpc.NewServer()

	pb.RegisterHelloServer(s, helloService)
	grpclog.Println("Listen on " + Address)
	s.Serve(listen)
}
