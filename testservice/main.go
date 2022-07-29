package main

import (
	"context"
	"log"
	"net"

	pb "github.com/EwanValentine/grpc-proxy-framework/testservice/gen/go/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello -
func (s *server) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

// CreateUser -
func (s *server) CreateUser(_ context.Context, in *pb.UserRequest) (*pb.UserResponse, error) {
	log.Println("user:", in.User)
	return &pb.UserResponse{Name: in.User.Name}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Panic(err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	reflection.Register(grpcServer)

	pb.RegisterGreeterServer(grpcServer, &server{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Panic(err)
	}
}
