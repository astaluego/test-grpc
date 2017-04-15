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

// New customer in database
func (s *server) New(ctx context.Context, in *customer.Customer) (*customer.Response, error) {

	// Verify if email is available
	var user customer.Customer
	has, err := engine.Where("email = ?", in.Email).Get(&user)
	if err != nil {
		return nil, err
	} else if has == true {
		return nil, fmt.Errorf("email %s already exists", in.Email)
	}

	// Insert customer on database
	_, err = engine.Insert(in)
	if err != nil {
		return nil, err
	}
	return &customer.Response{Message: "New " + in.Email}, nil
}

// Edit customer in database
func (s *server) Edit(ctx context.Context, in *customer.Customer) (*customer.Response, error) {
	return &customer.Response{Message: "Edit " + in.Email}, nil
}

// Delete customer in database
func (s *server) Delete(ctx context.Context, in *customer.Customer) (*customer.Response, error) {

	return &customer.Response{Message: "Delete " + in.Email}, nil
}

// List customer of database
func (s *server) List(ctx context.Context, in *customer.Customer) (*customer.Response, error) {
	var users []customer.Customer
	var results string

	//Select customer in database
	err := engine.Find(&users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		results = "empty"
	}
	for i, user := range users {
		if i != 0 {
			results += ", "
		}
		results += user.Email
	}
	return &customer.Response{Message: "List " + results}, nil
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
