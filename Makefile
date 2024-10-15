# Build the application
build:
	go build -o bin/api-gatekeeper ./cmd/api_gatekeeper

# Build and run the application
run: build
	./bin/api-gatekeeper -config=$(config)

# Run the application with the example/config.yaml file
run/dev:
	go run ./cmd/api_gatekeeper -config=./example/config.yaml
