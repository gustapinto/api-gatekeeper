build/linux:
	GOOS=linux GOARCH=amd64 go build -o bin/api-gatekeeper-linux-amd64 ./cmd/api_gatekeeper
	GOOS=linux GOARCH=arm64 go build -o bin/api-gatekeeper-linux-arm64 ./cmd/api_gatekeeper

build/windows:
	GOOS=windows GOARCH=amd64 go build -o bin/api-gatekeeper-windows-amd64.exe ./cmd/api_gatekeeper
	GOOS=windows GOARCH=arm64 go build -o bin/api-gatekeeper-windows-arm64.exe ./cmd/api_gatekeeper

build/macos:
	GOOS=darwin GOARCH=amd64 go build -o bin/api-gatekeeper-macos-amd64 ./cmd/api_gatekeeper
	GOOS=darwin GOARCH=arm64 go build -o bin/api-gatekeeper-macos-arm64 ./cmd/api_gatekeeper

build: build/linux build/windows build/macos

run:
	go run ./cmd/api_gatekeeper -config=./example/config.yaml
