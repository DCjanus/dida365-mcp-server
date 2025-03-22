package service

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	apiv1 "github.com/dcjanus/dida365-mcp-server/gen/proto/api/v1"
)

type Dida365MCP struct {
	logger *zap.Logger
	apiv1.UnimplementedData365MCPServer
}

func NewDida365MCP(logger *zap.Logger) *Dida365MCP {
	return &Dida365MCP{logger: logger}
}

func (d *Dida365MCP) Ping(ctx context.Context, req *emptypb.Empty) (*wrapperspb.StringValue, error) {
	return &wrapperspb.StringValue{Value: "Pong"}, nil
}
