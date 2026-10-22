package main

// go test -v -run TestAuthorizeTool
// go test -v  .

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestAuthorizeTool(t *testing.T) {
	validToken := signTestToken(t, tokenClaims{
		User: "alice",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, jwtSecret())
	expiredToken := signTestToken(t, tokenClaims{
		User: "alice",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}, jwtSecret())
	missingExpirationToken := signTestToken(t, tokenClaims{User: "alice"}, jwtSecret())
	wrongSecretToken := signTestToken(t, tokenClaims{
		User: "alice",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}, []byte("wrong secret"))

	tests := []struct {
		name    string
		header  string
		wantErr bool
	}{
		{name: "valid bearer token", header: "Bearer " + validToken},
		{name: "missing token", wantErr: true},
		{name: "wrong scheme", header: "Basic " + validToken, wantErr: true},
		{name: "expired token", header: "Bearer " + expiredToken, wantErr: true},
		{name: "missing expiration", header: "Bearer " + missingExpirationToken, wantErr: true},
		{name: "invalid signature", header: "Bearer " + wrongSecretToken, wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := &mcp.CallToolRequest{Extra: &mcp.RequestExtra{Header: make(http.Header)}}
			if test.header != "" {
				req.Extra.Header.Set("Authorization", test.header)
			}

			err := authorizeTool(req)
			if (err != nil) != test.wantErr {
				t.Fatalf("authorizeTool() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestProtectedToolsRejectMissingAuthorization(t *testing.T) {
	if _, _, err := add(nil, &mcp.CallToolRequest{}, addInput{A: 1, B: 2}); err == nil {
		t.Fatal("add() accepted a call without authorization")
	}

	before := counter.Load()
	if _, _, err := counterHandler(nil, &mcp.CallToolRequest{}, struct{}{}); err == nil {
		t.Fatal("counterHandler() accepted a call without authorization")
	}
	if got := counter.Load(); got != before {
		t.Fatalf("counter changed without authorization: got %d, want %d", got, before)
	}
}

func signTestToken(t *testing.T, claims tokenClaims, secret []byte) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("sign JWT: %v", err)
	}
	return signed
}
