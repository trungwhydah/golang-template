FROM golang:latest

WORKDIR /modules
RUN go install github.com/swaggo/swag/cmd/swag@latest
COPY go.mod go.sum /modules/
RUN go mod download

ENV PROJECT_DIR=/app \
    GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /app
RUN mkdir "/build"
COPY . .

RUN go get github.com/githubnemo/CompileDaemon
RUN go install github.com/githubnemo/CompileDaemon
RUN swag init -g internal/core_service/api/router.go -o internal/core_service/docs

ENTRYPOINT CompileDaemon -build="go build -o /build/app ./cmd/core_service" -command="/build/app"