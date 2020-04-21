package main

import (
	"context"
	"log"

	pb "github.com/sei-ri/grpc-gateway-rest-api-example/proto/greeter"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		log.Println(err)
	}

	client := pb.NewGreeterClient(conn)

	requestBody := &pb.HelloRequest{
		Name: "sushi",
	}
	response, err := client.SayHello(context.Background(), requestBody)
	if err != nil {
		log.Printf("SayHello error: %v", err)
	}
	log.Printf("SayHello response: %v", response.Message)


	payload := &pb.HelloRequest{
		Name: "ramen",
	}
	response, err = client.SayHelloAgain(context.Background(), payload)
	if err != nil {
		log.Printf("SayHelloAgain error: %v", err)
	}
	log.Printf("SayHelloAgain response: %v", response.Message)
}
