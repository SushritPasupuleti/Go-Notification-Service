include .env

run:
	@echo "Running server..."
	@go run main.go

build:
	@echo "Building server..."
	@GOOS=linux GOARCH=amd64 go build -o bin/server main.go 
	@GOOS=windows GOARCH=amd64 go build -o bin/server.exe main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/server.darwin main.go
	@echo "Done building for linux, windows and mac (x64 only)"
	@echo "Binaries are in bin/ directory"

run_build:
	@echo "Running server..."
	@./bin/server

test:
	@echo "Running tests..."
	@go test -v ./...

lint:
	@echo "Running linter..."
	@golangci-lint run

benchmark:
	@echo "Running benchmark..."
	@wrk2 -t2 -c100 -d30s -R2000 -s benchmark.lua http://localhost:${PORT}/producer/send 
	# For 30s, 2000 requests per second, 100 connections, 2 threads

kafka_up:
	@echo "Starting kafka..."
	@sudo docker compose up -d

kafka_down:
	@echo "Stopping kafka..."
	@sudo docker compose down

kafka_stop:
	@echo "Stopping kafka..."
	@sudo docker compose stop

kafka_rm:
	@echo "Removing kafka..."
	@sudo docker compose rm

dev:
	@air

swagger_docgen:
	@echo "Generating Swagger Docs..."
	@echo ""
	@rm -rf docs/
	@swag init
	@echo ""
	@echo "Done generating docs."
	@echo ""
	@echo "Visit: http://localhost:${PORT}/swagger/index.html"

docgen: swagger_docgen
	@echo ""
	@echo "Generating OpenAPIv3 Docs..."
	@echo ""
	@rm -rf docs/openapi.yaml
	@npx -p swagger2openapi swagger2openapi --yaml --outfile docs/openapi.yaml "http://localhost:${PORT}/swagger/doc.json"
	@echo "Done generating OpenAPIv3 Docs."
