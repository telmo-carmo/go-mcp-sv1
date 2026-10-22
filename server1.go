package main

/*
run:

npx @modelcontextprotocol/inspector

*/

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	MY_JWT_SECRET = "A Very Secret Key for JWT!"
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

type tokenOutput struct {
	Token string `json:"value" jsonschema:"The JWT token value"`
}

type tokenInput struct {
	Username string `json:"username" jsonschema:"The username to include in the JWT"`
}

func tokenHandler(_ context.Context, _ *mcp.CallToolRequest, req tokenInput) (*mcp.CallToolResult, any, error) {
	claims := jwt.MapClaims{
		"user": req.Username,
		"exp":  time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	key := MY_JWT_SECRET
	if envKey := os.Getenv("MY_JWT_SECRET"); envKey != "" {
		key = envKey
	}
	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return nil, nil, err
	}
	log.Printf("[INFO] Generated JWT for user %s: %s", req.Username, signedToken)
	return nil, tokenOutput{Token: signedToken}, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

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
	mcp.AddTool(server, &mcp.Tool{
		Name:        "token",
		Description: "Generate a new JWT.",
	}, tokenHandler)

	mux := http.NewServeMux()
	mux.Handle("/mcp", mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil))

	const address = "localhost:8080"
	log.Printf("[INFO] MCP server listening at http://%s/mcp", address)
	log.Fatal(http.ListenAndServe(address, mux))
}
