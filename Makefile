
# Go parameters
GCP_PROJECT_NAME=s9-demo
BINARY_NAME=tfeel

all: test
build:
	go build -o ./bin/$(BINARY_NAME) -v

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/$(BINARY_NAME)

test:
	go test -v ./...

clean:
	go clean
	rm -f ./bin/$(BINARY_NAME)

run: build
	bin/$(BINARY_NAME)

deps:
	go get github.com/golang/dep/cmd/dep
	dep ensure