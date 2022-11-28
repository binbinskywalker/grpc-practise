package main

import (
	"context"
	"fmt"
	"net"

	pb "my_grpc/proto/hello"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

const (
	Address = "127.0.0.1:8888"
)

type helloService struct{}

var HelloService = helloService{}

func (h helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponce, error) {
	resp := new(pb.HelloResponce)
	resp.Message = fmt.Sprintf("Hello %s,", in.Name)
	return resp, nil
}
func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("fail to listen :%v", err)
	}
	s := grpc.NewServer()

	pb.RegisterHelloServer(s, HelloService)
	grpclog.Println("Listen on " + Address)
	s.Serve(listen)
}
