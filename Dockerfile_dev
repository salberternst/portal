FROM golang:1.21.1 as base

RUN curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.42.0

WORKDIR /app

CMD air --build.cmd "go build -o bin/server cmd/server/main.go" --build.bin "./bin/server"