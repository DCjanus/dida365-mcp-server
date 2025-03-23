package main

import (
	"flag"
	"os"

	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

func main() {
	accessToken := os.Getenv("MCP_ACCESS_TOKEN")
	flag.StringVar(&accessToken, "access_token", "", "The access token to use for the MCP server, can be set using the MCP_ACCESS_TOKEN environment variable")
	flag.Parse()

	if accessToken == "" {
		flag.Usage()
		os.Exit(1)
	}

	done := make(chan struct{})
	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())
	err := server.Serve()
	if err != nil {
		panic(err)
	}
	<-done
}
