package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"path"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	pb "github.com/sei-ri/grpc-gateway-rest-api-example/proto/greeter"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"

	"github.com/elazarl/go-bindata-assetfs"
	swagger "github.com/sei-ri/grpc-gateway-rest-api-example/web/ui/swagger/data"
)

var (
	ServerAddr string
	SwaggerDir string
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
	mux.HandleFunc("/swagger/", swaggerHandle(SwaggerDir))
	serveSwaggerUI(mux)

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

// h2cHandleFunc returns switch http handler simulation http2
func h2cHandleFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

// swaggerHandle returns swagger specification files located under "/swagger/"
func swaggerHandle(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "swagger.json") {
			log.Printf("Not Found: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}

		log.Printf("Serving %s", r.URL.Path)
		p := strings.TrimPrefix(r.URL.Path, "/swagger/")
		p = path.Join(dir, p)

		log.Printf("Serving swagger-file: %s", p)

		http.ServeFile(w, r, p)
	}
}

// frontend swagger-ui
func serveSwaggerUI(mux *http.ServeMux) {
	fileServer := http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
		Prefix:   "third_party/swagger-ui",
	})
	prefix := "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, fileServer))
}