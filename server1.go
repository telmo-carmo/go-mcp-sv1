package main

/*
run:

npx @modelcontextprotocol/inspector

*/

import (
	"context"
	"log"
	"net/http"
	"sync/atomic"

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
	A float64 `json:"a" jsonschema:"The first number to add"`
	B float64 `json:"b" jsonschema:"The second number to add"`
}
type addOutput struct {
	Result float64 `json:"result" jsonschema:"The calculated sum of the two numbers"`
}

func add(_ context.Context, _ *mcp.CallToolRequest, input addInput) (*mcp.CallToolResult, any, error) {
	sum := input.A + input.B
	log.Printf("[INFO] Adding %.2f and %.2f, result: %.2f", input.A, input.B, sum)
	return nil, addOutput{Result: sum}, nil
}

var counter atomic.Int64

type counterOutput struct {
	Value int64 `json:"value" jsonschema:"The counter value"`
}

func counterHandler(_ context.Context, _ *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	nv := counter.Add(1)
	log.Printf("[INFO] Counter incremented, new value: %d", nv)
	return nil, counterOutput{Value: nv}, nil
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
		Description: "Add two numbers.",
	}, add)
	mcp.AddTool(server, &mcp.Tool{
		Name:        "counter",
		Description: "Increment and return the counter.",
	}, counterHandler)

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil))

	const address = "localhost:8080"
	log.Printf("[INFO] MCP server listening at http://%s/mcp", address)
	log.Fatal(http.ListenAndServe(address, mux))
}
