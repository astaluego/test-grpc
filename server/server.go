package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	service "github.com/astaluego/test-grpc/server/pkg/protobuf"
	"github.com/astaluego/test-grpc/server/pkg/protobuf/customer"
	"github.com/go-xorm/xorm"
	_ "github.com/lib/pq"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	port = ":4242"
)

var engine = InitEngine()

type server struct{}

//TODO verifier si l'adresse n'existe pas
func (s *server) New(ctx context.Context, in *customer.Customer) (*customer.Response, error) {
	_, err := engine.Insert(in)
	if err != nil {
		return nil, err
	}
	results, err := engine.Query("select * from customer")
	if err != nil {
		return nil, err
	}
	fmt.Print(results)
	return &customer.Response{Message: "New " + in.Email}, nil
}

func (s *server) Edit(ctx context.Context, in *customer.Customer) (*customer.Response, error) {
	return &customer.Response{Message: "Edit " + in.Email}, nil
}

func (s *server) Delete(ctx context.Context, in *customer.Customer) (*customer.Response, error) {
	return &customer.Response{Message: "Delete " + in.Email}, nil
}

func (s *server) List(ctx context.Context, in *customer.Customer) (*customer.Response, error) {
	return &customer.Response{Message: "List " + in.Email}, nil
}

// InitEngine init database
func InitEngine() *xorm.Engine {
	database, err := xorm.NewEngine("postgres", "host=127.0.0.1 port=5432 user=turlutte password=anastasia dbname=turlutte sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	err = database.CreateTables(&customer.Customer{})
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}
	return database
}

func main() {
	fmt.Println("Server")
	flag.Parse()

	// engine
	engine.ShowSQL(true)

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	service.RegisterRouteServer(s, &server{})
	s.Serve(lis)
}
