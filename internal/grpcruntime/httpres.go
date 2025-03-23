package grpcruntime

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/proto"

	"github.com/dcjanus/dida365-mcp-server/gen/model"
)

func TemporaryRedirectForwardResponseOption() runtime.ServeMuxOption {
	return runtime.WithForwardResponseOption(func(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
		if rr, ok := message.(*model.TemporaryRedirectResponse); ok {
			writer.Header().Set("Location", rr.Location)
			writer.WriteHeader(http.StatusTemporaryRedirect)
			return nil
		}
		return nil
	})
}
func HTMLResponseForwardResponseOption() runtime.ServeMuxOption {
	return runtime.WithForwardResponseOption(func(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
		if rr, ok := message.(*model.HTMLResponse); ok {
			writer.Header().Set("Content-Type", "text/html")
			writer.WriteHeader(http.StatusOK)
			_, err := writer.Write([]byte(rr.Html))
			return err
		}
		return nil
	})
}
