# syntax=docker/dockerfile:1

ARG GO_VERSION="1.24-alpine"

FROM golang:${GO_VERSION} AS build

WORKDIR /app

RUN apk add --no-cache --update build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ./build/server ./cmd/server

RUN go test -v ./...

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=build /app/build/server /app/amazing-core

CMD ["/app/amazing-core"]
