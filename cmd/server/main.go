package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	helloworldv1 "github.com/DCjanus/dida365-mcp-server/gen/proto/helloworld/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type server struct {
	helloworldv1.UnimplementedHelloServiceServer
}

// SayHello 实现 HelloService 服务
func (s *server) SayHello(ctx context.Context, req *helloworldv1.SayHelloRequest) (*helloworldv1.SayHelloResponse, error) {
	return &helloworldv1.SayHelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}

func main() {
	// 创建 TCP 监听器
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 创建 gRPC 服务器
	s := grpc.NewServer()
	helloworldv1.RegisterHelloServiceServer(s, &server{})

	// 创建 gRPC-Gateway mux
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 注册 gRPC-Gateway handler
	if err := helloworldv1.RegisterHelloServiceHandlerFromEndpoint(
		context.Background(),
		gwmux,
		"localhost:8080",
		opts,
	); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	// 创建支持 h2c 的 HTTP 服务器
	h2s := &http2.Server{
		MaxReadFrameSize:     1048576, // 1MB
		MaxConcurrentStreams: 250,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && r.Header.Get("Content-Type") == "application/grpc" {
			s.ServeHTTP(w, r)
			return
		}
		gwmux.ServeHTTP(w, r)
	})

	// 使用 h2c 包装 handler
	h2cHandler := h2c.NewHandler(handler, h2s)

	httpServer := &http.Server{
		Handler:           h2cHandler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	// 启动服务器
	log.Printf("Server listening on %s", lis.Addr().String())
	if err := httpServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
