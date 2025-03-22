package service

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/dcjanus/dida365-mcp-server/gen/api"
)

type Dida365 struct {
	logger *zap.Logger
	api.UnimplementedData365ServiceServer
}

func NewDida365MCP(logger *zap.Logger) *Dida365 {
	return &Dida365{logger: logger}
}

func (d *Dida365) Ping(ctx context.Context, req *emptypb.Empty) (*wrapperspb.StringValue, error) {
	return &wrapperspb.StringValue{Value: "Pong"}, nil
}
