package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	service "github.com/astaluego/test-grpc/client/pkg/protobuf"
	"github.com/astaluego/test-grpc/client/pkg/protobuf/user"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "testdata/ca.pem", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:4242", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
	email              = flag.String("email", "user@domain.com", "Email adress")
)

func main() {
	fmt.Println("Client")
	flag.Parse()

	// Init connection
	conn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Can't connect to server: %v", err)
	}
	defer conn.Close()

	// Init service Route
	client := service.NewRouteClient(conn)

	var response *user.Response

	// Read
LOOP:
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		switch scanner.Text() {
		case "1":
			response, err = client.New(context.Background(), &user.User{Email: *email})
			break
		case "2":
			response, err = client.Edit(context.Background(), &user.User{Email: *email})
			break
		case "3":
			response, err = client.Delete(context.Background(), &user.User{Email: *email})
			break
		case "4":
			response, err = client.List(context.Background(), &user.User{Email: *email})
			break
		case "exit":
			return
		default:
			fmt.Println("Command not found")
			goto LOOP
		}
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Succes: %s", response.Message)
	}
}
