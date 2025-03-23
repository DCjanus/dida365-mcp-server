package main

import (
	"context"
	"flag"
	"os"

	"github.com/bufbuild/protovalidate-go"
	"github.com/cockroachdb/errors"
	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"

	"github.com/dcjanus/dida365-mcp-server/gen/conf"
	"github.com/dcjanus/dida365-mcp-server/internal/utils"
)

func main() {
	utils.LoadDotEnvs()

	verbose := false
	accessToken := ""
	flag.StringVar(&accessToken, "access_token", "", "The access token to use for the MCP server, can be set using the MCP_ACCESS_TOKEN environment variable")
	flag.BoolVar(&verbose, "verbose", false, "Whether to enable verbose logging")
	flag.Parse()

	if accessToken == "" {
		accessToken = os.Getenv("MCP_ACCESS_TOKEN")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loggingConfig := &conf.Logging{Level: "info"}
	if verbose {
		loggingConfig.Level = "debug"
	}

	if err := protovalidate.Validate(loggingConfig); err != nil {
		panic(errors.Wrap(err, "invalid logging config"))
	}

	log, err := utils.NewLogger(loggingConfig)
	if err != nil {
		panic(errors.Wrap(err, "failed to create logger"))
	}
	defer func() { _ = log.Sync() }()

	if accessToken == "" {
		flag.Usage()
		log.Fatal("access token is required")
		os.Exit(1)
	}

	srv := server.NewMCPServer(
		"Dida365 MCP Server",
		"0.1.0",
	)

	wrapper, err := NewDidaWrapper(ctx, log, accessToken)
	if err != nil {
		log.Fatal("failed to create dida tools", zap.Error(err))
	}
	srv.AddTools(wrapper.Tools()...)

	log.Info("starting MCP server")
	if err := server.ServeStdio(srv); err != nil {
		log.Fatal("failed to serve MCP server", zap.Error(err))
	}
}
