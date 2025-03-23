package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/bufbuild/protovalidate-go"
	"github.com/cockroachdb/errors"
	mcpgolang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"go.uber.org/zap"

	"github.com/dcjanus/dida365-mcp-server/gen/conf"
	"github.com/dcjanus/dida365-mcp-server/internal/utils"
)

func main() {
	verbose := false
	accessToken := os.Getenv("MCP_ACCESS_TOKEN")
	flag.StringVar(&accessToken, "access_token", "", "The access token to use for the MCP server, can be set using the MCP_ACCESS_TOKEN environment variable")
	flag.BoolVar(&verbose, "verbose", false, "Whether to enable verbose logging")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	loggingConfig := &conf.Logging{
		Level: "info",
	}
	if verbose {
		loggingConfig.Level = "debug"
	}

	if err := protovalidate.Validate(loggingConfig); err != nil {
		panic(errors.Wrap(err, "invalid logging config"))
	}

	if accessToken == "" {
		flag.Usage()
		os.Exit(1)
	}
	log, err := utils.NewLogger(loggingConfig)
	if err != nil {
		panic(errors.Wrap(err, "failed to create logger"))
	}
	defer func() {
		_ = log.Sync()
	}()

	tools, err := NewDidaTools(ctx, log, accessToken)
	if err != nil {
		log.Fatal("failed to create dida tools", zap.Error(err))
	}

	server := mcpgolang.NewServer(stdio.NewStdioServerTransport())
	if err := tools.Register(server); err != nil {
		log.Fatal("failed to register dida tools", zap.Error(err))
	}

	if err := server.Serve(); err != nil {
		log.Fatal("failed to serve MCP server", zap.Error(err))
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	log.Info("received signal, shutting down")
}
