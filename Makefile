
go-mcp-sv1: server1.go
	@echo "Building go-mcp-sv1..."
	@go build -o go-mcp-sv1 server1.go

upd:
	go get -u -v
	go mod tidy
