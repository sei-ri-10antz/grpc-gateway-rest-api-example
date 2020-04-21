package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/sei-ri/grpc-gateway-rest-api-example/proto/greeter"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

var (
	ServerAddr string
)

func Run() (err error) {
	// New listen
	conn, err := net.Listen("tcp", ServerAddr)
	if err != nil {
		log.Printf("Failed to TCP listen: %v", err)
	}

	// New http server
	server := newHttpServer()
	
	log.Printf("gRPC server listening at: %v", conn.Addr())

	// Serve server
	if err := server.Serve(conn); err != nil {
		log.Printf("Failed to server %v", err)
	}

	return nil
}

func newHttpServer() *http.Server {
	g := newGrpc()
	gw, err := newGateway()
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", gw)

	return &http.Server{
		Addr: ServerAddr,
		Handler: h2cHandleFunc(g, mux), // switch handler
	}
}

func newGrpc() *grpc.Server {
	var opts []grpc.ServerOption

	// New gRPC server
	server := grpc.NewServer(opts...)

	// Register gRPC service pb
	pb.RegisterGreeterServer(server, NewGreeterServerHandler())

	return server
}

func newGateway() (http.Handler, error) {
	ctx := context.Background()

	opts := []grpc.DialOption{grpc.WithInsecure()}

	// New gw server
	gwMux := runtime.NewServeMux()

	// Register gateway endpoint
	if err := pb.RegisterGreeterHandlerFromEndpoint(ctx, gwMux, ServerAddr, opts); err != nil {
		return nil, err
	}

	return gwMux, nil
}

// この関数は、requestリクエストがrpc clientで開始されるか、またはrest apiで開始されるかを判別するために使用されます。
// 異なるrequestに応じてServeHTTPサービスを登録して処理する。
// r.ProtoMajor==2は、HTTP/2を代表する。
func h2cHandleFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}