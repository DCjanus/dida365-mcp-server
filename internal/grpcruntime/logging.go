package grpcruntime

import (
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(map[string]string{})
		}
		fields := []zap.Field{
			zap.String("grpc_method", info.FullMethod),
			zap.Any("request", req),
		}

		begin := time.Now()
		reply, err := handler(ctx, req)

		fields = append(fields, zap.Duration("duration", time.Since(begin)))

		if st, ok := status.FromError(err); ok {
			fields = append(fields, zap.String("code", st.Code().String()))
		}

		if httpPath := md.Get("http_path"); len(httpPath) > 0 && httpPath[0] != "" {
			fields = append(fields, zap.String("http_path", httpPath[0]))
		}
		if httpHost := md.Get("http_host"); len(httpHost) > 0 && httpHost[0] != "" {
			fields = append(fields, zap.String("http_host", httpHost[0]))
		}
		if httpMethod := md.Get("http_method"); len(httpMethod) > 0 && httpMethod[0] != "" {
			fields = append(fields, zap.String("http_method", httpMethod[0]))
		}

		log.Info("request", fields...)
		return reply, err
	}
}

func WithHTTPMetadata(log *zap.Logger) runtime.ServeMuxOption {
	return runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
		return metadata.New(map[string]string{
			"http_path":   req.URL.Path,
			"http_host":   req.Host,
			"http_method": req.Method,
		})
	})
}
