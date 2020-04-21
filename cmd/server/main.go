package main

import (
	"github.com/sei-ri/grpc-gateway-rest-api-example/server"
)

func main() {
	server.ServerAddr = ":8080"
	server.Run()
}
