
go-mcp-sv1: server1.go
	@echo "Building go-mcp-sv1..."
	@go build -o go-mcp-sv1 server1.go

go-mcp-sv1.exe: server1.go
	go build -o go-mcp-sv1.exe server1.go
	
	
upd:
	go get -u -v
	go mod tidy

clean:
	@echo "Cleaning up..."
	@rm -f go-mcp-sv1 go-mcp-sv1.exe
	