package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/cockroachdb/errors"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/dcjanus/dida365-mcp-server/gen/api"
	"github.com/dcjanus/dida365-mcp-server/internal/grpcruntime"
	"github.com/dcjanus/dida365-mcp-server/internal/service"
	"github.com/dcjanus/dida365-mcp-server/internal/utils"
)

func main() {
	configPath := flag.String("config", "/etc/dida365-oauth-server/config.yaml", "path to config file")
	version := flag.Bool("version", false, "print version")
	flag.Parse()

	if *version {
		fmt.Println(utils.Version)
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	utils.LoadDotEnvs()

	cfg, err := utils.LoadConfig(*configPath)
	if err != nil {
		panic(errors.Wrap(err, "failed to load config"))
	}

	log, err := utils.NewLogger(cfg.GetLogging())
	if err != nil {
		panic(errors.Wrap(err, "failed to create logger"))
	}

	lis, err := net.Listen("tcp", cfg.GetServer().GetListen())
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(
		grpcmiddleware.ChainUnaryServer(
			grpcruntime.ValidateInterceptor(),
			grpcruntime.LoggingInterceptor(log),
		),
	))
	api.RegisterDida365OAuthServiceServer(srv, service.NewDida365AuthService(log, cfg))

	mux := runtime.NewServeMux(
		grpcruntime.TemporaryRedirectForwardResponseOption(),
		grpcruntime.WithHTTPMetadata(log),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:     true,
					EmitDefaultValues: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)
	if err := mux.HandlePath(http.MethodGet, "/metrics", func(w http.ResponseWriter, r *http.Request, _ map[string]string) {
		promhttp.Handler().ServeHTTP(w, r)
	}); err != nil {
		log.Fatal("failed to register metrics handler", zap.Error(err))
	}
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	if err := api.RegisterDida365OAuthServiceHandlerFromEndpoint(ctx, mux, lis.Addr().String(), opts); err != nil {
		log.Fatal("failed to register gRPC gateway", zap.Error(err))
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

	log.Info("Server starting", zap.String("listen", lis.Addr().String()))
	if err := httpServer.Serve(lis); err != nil {
		log.Fatal("failed to serve", zap.Error(err))
	}
}
