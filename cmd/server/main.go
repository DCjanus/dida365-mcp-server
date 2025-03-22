package main

import (
	"context"
	"flag"
	"net"
	"net/http"
	"time"

	"github.com/cockroachdb/errors"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	apiv1 "github.com/dcjanus/dida365-mcp-server/gen/proto/api/v1"
	"github.com/dcjanus/dida365-mcp-server/internal/middleware"
	"github.com/dcjanus/dida365-mcp-server/internal/service"
	"github.com/dcjanus/dida365-mcp-server/internal/utils"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := utils.LoadConfig(*configPath)
	if err != nil {
		panic(errors.Wrap(err, "failed to load config"))
	}

	logger, err := utils.NewLogger(cfg.GetLogging())
	if err != nil {
		panic(errors.Wrap(err, "failed to create logger"))
	}

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Fatal("failed to listen", zap.Error(err))
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcmiddleware.ChainUnaryServer(
			middleware.Validate(),
		),
	))
	apiv1.RegisterData365MCPServer(srv, service.NewDida365MCP(logger))

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:     true,
			EmitDefaultValues: true,
		},
	}))
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := apiv1.RegisterData365MCPHandlerFromEndpoint(ctx, mux, "localhost:8080", opts); err != nil {
		logger.Fatal("failed to register gRPC gateway", zap.Error(err))
	}

	h2s := &http2.Server{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && r.Header.Get("Content-Type") == "application/grpc" {
			srv.ServeHTTP(w, r)
			return
		}
		mux.ServeHTTP(w, r)
	})

	h2cHandler := h2c.NewHandler(handler, h2s)

	httpServer := &http.Server{
		Handler:           h2cHandler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	logger.Info("Server starting", zap.String("listen", lis.Addr().String()))
	if err := httpServer.Serve(lis); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
