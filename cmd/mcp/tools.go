package main

import (
	"context"
	"encoding/json"

	"github.com/cockroachdb/errors"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"go.uber.org/zap"

	"github.com/dcjanus/dida365-mcp-server/internal/dida"
)

type DidaTools struct {
	log *zap.Logger
	cli *dida.Client
}

func NewDidaTools(ctx context.Context, log *zap.Logger, token string) (*DidaTools, error) {
	cli := dida.NewClient(log, token)
	if _, err := cli.ListProjects(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to check dida token")
	}

	return &DidaTools{
		log: log.With(zap.String("component", "mcp.DidaTools")),
		cli: cli,
	}, nil
}

func (t *DidaTools) Register(server *mcpgolang.Server) error {
	if err := server.RegisterTool("dida.ListProjects", "List all projects", func(ctx context.Context) (*mcpgolang.ToolResponse, error) {
		res, err := t.cli.ListProjects(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to list projects")
		}
		inJson, err := json.Marshal(res)
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal projects")
		}
		content := mcpgolang.NewTextContent(string(inJson))
		return mcpgolang.NewToolResponse(content), nil
	}); err != nil {
		return errors.Wrap(err, "failed to register dida.ListProjects")
	}

	return nil
}
