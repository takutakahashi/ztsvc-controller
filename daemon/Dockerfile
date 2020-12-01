FROM golang as build

WORKDIR /go/src
COPY go.mod /go/src/
COPY go.sum /go/src/

RUN go mod download

COPY . /go/src/
RUN make build

FROM ubuntu
RUN apt update && apt install -y curl gpg
RUN curl -s https://install.zerotier.com | bash
COPY --from=build /go/src/dist/daemon /bin/ztdaemon
ENTRYPOINT ["/bin/ztdaemon"]
