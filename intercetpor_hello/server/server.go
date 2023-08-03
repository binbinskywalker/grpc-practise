package main

import (
	"fmt"
	"net"

	pb "my_grpc/proto/hello" // 引入proto包

	"cloud.google.com/go/functions/metadata"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
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

	var opts []grpc.ServerOption

	creds, err := credentials.NewServerTLSFromFile("../keys/server.pem", "../keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to genenrate credentials:%v", err)
	}
	opts = append(opts, grpc.Creds(creds))
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	s := grpc.NewServer(opts...)

	pb.RegisterHelloServer(s, HelloService)
	grpclog.Println("Listen on " + Address + "with TLS +Token + Interceptor")
	s.Serve(listen)
}

func auth(ctx context.Context) error {
	md, ok := metadata.FromContext(ctx)
	if ok != nil {
		return grpc.Errorf(codes.Unauthenticated, "no Token auth infos")
	}
	var (
		appid  string
		appkey string
	)

	if val, ok := md["appid"]; ok {
		appid = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appkey = val[1]
	}
	if appid != "101010" || appkey != "i am key" {
		return grpc.Errorf(codes.Unauthenticated, "token is invalid:appid:%s, appkey:%s", appid, appkey)
	}
	return nil
}

func interceptor(ctx context.Context, req interface{}, infos *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)

}
