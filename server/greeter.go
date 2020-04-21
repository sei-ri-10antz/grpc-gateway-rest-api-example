package server

import (
	"context"
	"log"

	pb "github.com/sei-ri/grpc-gateway-rest-api-example/proto/greeter"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type greeterServerHandler struct{}

func NewGreeterServerHandler() pb.GreeterServer {
	return new(greeterServerHandler)
}

func (self *greeterServerHandler) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("recevied Greeter.SayHello %v", req)

	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid name")
	}

	return &pb.HelloReply{Message: "Say Hello " + req.Name}, nil
}

func (self *greeterServerHandler) SayHelloAgain(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("recevied Greeter.SayHelloAgain %v", req)

	if len(req.Name) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid name")
	}

	resp := new(pb.HelloReply)
	resp.Message = "Say Hello Again " + req.Name
	return resp, nil
}
