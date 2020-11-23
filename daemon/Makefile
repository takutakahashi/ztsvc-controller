all: build test

test:
	dist/daemon --config example/config.yaml

build:
	GO111MODULE=on go build -o dist/daemon cmd/cmd.go
