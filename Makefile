.PHONY: all test build clean

all: clean test build

build: 
	mkdir -p build
	go build -o build ./...

test:
	go test -v -coverprofile=tests/results/cover.out ./...

verify:
	golangci-lint run -c .golangci.yaml --deadline=5m

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

clean:
	rm -rf build/*
	go clean ./...

container:
	podman build -t  quay.io/luzuccar/lab-tekton-emulator-cicd:v0.0.1 .

push:
	podman push quay.io/luigizuccarelli/lab-tekton-emulator-cicd:v0.0.1
