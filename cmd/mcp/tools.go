package main

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"

	"github.com/dcjanus/dida365-mcp-server/internal/dida"
)

type DidaWrapper struct {
	log *zap.Logger
	cli *dida.Client
	ctx context.Context
}

func NewDidaWrapper(ctx context.Context, log *zap.Logger, token string) (*DidaWrapper, error) {
	cli := dida.NewClient(log, token)
	if _, err := cli.ListProjects(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to check dida token")
	}

	return &DidaWrapper{
		log: log.With(zap.String("component", "mcp.DidaWrapper")),
		cli: cli,
		ctx: ctx,
	}, nil
}

func (t *DidaWrapper) Tools() []server.ServerTool {
	return []server.ServerTool{
		t.ListProjects(t.ctx),
	}
}

func (t *DidaWrapper) ListProjects(ctx context.Context) server.ServerTool {
	return server.ServerTool{
		Tool: mcp.NewTool("list_projects", mcp.WithDescription("List all projects, projects are the top level container for tasks")),
		Handler: func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			projects, err := t.cli.ListProjects(ctx)
			if err != nil {
				return nil, errors.Wrap(err, "failed to list projects")
			}
			inJson, err := json.Marshal(projects)
			if err != nil {
				return nil, errors.Wrap(err, "failed to marshal projects")
			}
			response := mcp.NewToolResultText(string(inJson))
			return response, nil
		},
	}
}
