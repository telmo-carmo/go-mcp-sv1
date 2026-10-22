package main

/*
run:

npx @modelcontextprotocol/inspector

*/

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const (
	MY_JWT_SECRET = "A Very Secret Key for JWT!"
)

type tokenClaims struct {
	User string `json:"user"`
	jwt.RegisteredClaims
}

func jwtSecret() []byte {
	secret := os.Getenv("MY_JWT_SECRET")
	if secret == "" {
		secret = MY_JWT_SECRET
	}
	return []byte(secret)
}

func authorizeTool(req *mcp.CallToolRequest) error {
	if req == nil || req.Extra == nil {
		return errors.New("Authorization bearer token is required")
	}
	fields := strings.Fields(req.Extra.Header.Get("Authorization"))
	if len(fields) != 2 || !strings.EqualFold(fields[0], "Bearer") {
		return errors.New("Authorization header must contain a bearer token")
	}

	claims := &tokenClaims{}
	token, err := jwt.ParseWithClaims(fields[1], claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected JWT signing method %q", token.Method.Alg())
		}
		return jwtSecret(), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return fmt.Errorf("invalid bearer JWT: %w", err)
	}
	if !token.Valid || claims.ExpiresAt == nil || strings.TrimSpace(claims.User) == "" {
		return errors.New("bearer JWT must be valid, unexpired, and contain user and exp claims")
	}
	return nil
}

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

func add(_ context.Context, req *mcp.CallToolRequest, input addInput) (*mcp.CallToolResult, any, error) {
	if err := authorizeTool(req); err != nil {
		return nil, nil, err
	}
	sum := input.A + input.B
	log.Printf("[INFO] Adding %.2f and %.2f, result: %.2f", input.A, input.B, sum)
	return nil, addOutput{Result: sum}, nil
}

var counter atomic.Int64

type counterOutput struct {
	Value int64 `json:"value" jsonschema:"The counter value"`
}

func counterHandler(_ context.Context, req *mcp.CallToolRequest, _ struct{}) (*mcp.CallToolResult, any, error) {
	if err := authorizeTool(req); err != nil {
		return nil, nil, err
	}
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
	signedToken, err := token.SignedString(jwtSecret())
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
