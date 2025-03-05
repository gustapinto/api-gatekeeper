build/linux:
	GOOS=linux GOARCH=amd64 go build -o bin/api-gatekeeper-rest-linux-amd64 ./cmd/api_gatekeeper_rest
	GOOS=linux GOARCH=arm64 go build -o bin/api-gatekeeper-rest-linux-arm64 ./cmd/api_gatekeeper_rest

build/windows:
	GOOS=windows GOARCH=amd64 go build -o bin/api-gatekeeper-rest-windows-amd64.exe ./cmd/api_gatekeeper_rest
	GOOS=windows GOARCH=arm64 go build -o bin/api-gatekeeper-rest-windows-arm64.exe ./cmd/api_gatekeeper_rest

build/macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/api-gatekeeper-rest-macos-amd64 ./cmd/api_gatekeeper_rest
	GOOS=darwin GOARCH=arm64 go build -o bin/api-gatekeeper-rest-macos-arm64 ./cmd/api_gatekeeper_rest

build: build/linux build/windows build/macos

run:
	go run ./cmd/api_gatekeeper_rest -config=./example/config.yaml
