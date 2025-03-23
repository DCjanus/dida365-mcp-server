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
