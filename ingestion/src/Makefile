# Define the binary name
BINARY_NAME=bootstrap

# Default target for building the binary
build:
	echo "Building the binary..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_NAME) main.go

# Custom target for SAM CLI (specific to the function name)
build-DataIngestionFunction: build

# Clean up generated files
clean:
	rm -f $(BINARY_NAME)

# Optional target to run tests
test:
	go test ./...

