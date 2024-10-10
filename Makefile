# Build the application
build:
	go build -o bin/api-gatekeeper ./cmd/api_gatekeeper

# Run the application with the example/config.yaml file
run:
	go run ./cmd/api_gatekeeper -config=./example/config.yaml
