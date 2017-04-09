package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	service "github.com/astaluego/test-grpc/server/pkg/protobuf"
	"github.com/astaluego/test-grpc/server/pkg/protobuf/user"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":4242"
)

func (s *server) New(ctx context.Context, in *user.User) (*user.Response, error) {
	return &user.Response{Message: "New " + in.Email}, nil
}

func (s *server) Edit(ctx context.Context, in *user.User) (*user.Response, error) {
	return &user.Response{Message: "Edit " + in.Email}, nil
}

func (s *server) Delete(ctx context.Context, in *user.User) (*user.Response, error) {
	return &user.Response{Message: "Delete " + in.Email}, nil
}

func (s *server) List(ctx context.Context, in *user.User) (*user.Response, error) {
	return &user.Response{Message: "List " + in.Email}, nil
}

type server struct{}

func main() {
	fmt.Println("Server")
	flag.Parse()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service.RegisterRouteServer(s, &server{})
	s.Serve(lis)
}
