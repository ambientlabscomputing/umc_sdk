.PHONY: all build test proto lint clean install

all: build

build:
	@echo "Building UMC SDK..."
	go build ./...

test:
	@echo "Running tests..."
	go test -v -race ./...

proto:
	@echo "Generating protobuf code..."
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/ua_kernel/v1/syscall.proto

lint:
	@echo "Running linters..."
	go vet ./...
	go fmt ./...

clean:
	@echo "Cleaning build artifacts..."
	go clean
	rm -f proto/ua_kernel/v1/*.pb.go

install:
	@echo "Installing as local module..."
	go mod tidy
	go mod download
