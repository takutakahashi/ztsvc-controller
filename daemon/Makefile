all: build test

test:
	go test ./... -coverprofile cover.out

build:
	GO111MODULE=on go build -o dist/daemon cmd/cmd.go
run:
	dist/daemon --token ${ZT_TOKEN} --networkID ${NETWORK_ID} --name test
