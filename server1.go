package main

import (
	"context"
	"log"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type greetInput struct {
	In string `json:"in"`
}

func greet(_ context.Context, _ *mcp.CallToolRequest, input greetInput) (*mcp.CallToolResult, any, error) {
	log.Printf("[WARN] Greeting %s", input.In)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Hello, " + input.In + "!"},
		},
	}, nil, nil
}

type addInput struct {
	A int `json:"a"`
	B int `json:"b"`
}
type addOutput struct {
	Result int `json:"result" jsonschema:"The calculated sum of the two numbers"`
}

func add(_ context.Context, _ *mcp.CallToolRequest, input addInput) (*mcp.CallToolResult, any, error) {
	sum := input.A + input.B
	log.Printf("[INFO] Adding %d and %d, result: %d", input.A, input.B, sum)
	return nil, addOutput{Result: sum}, nil
}

func main() {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "go-mcp-sv1",
		Version: "1.0.0",
	}, nil)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "greet",
		Description: "Greet someone by name.",
	}, greet)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "add",
		Description: "Add two int numbers.",
	}, add)

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil))

	const address = "localhost:8080"
	log.Printf("[INFO] MCP server listening at http://%s/mcp", address)
	log.Fatal(http.ListenAndServe(address, mux))
}
