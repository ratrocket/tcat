test:
	go test ./...

build:
	go build

install:
	go install

cov:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
