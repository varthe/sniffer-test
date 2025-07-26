# syntax=docker/dockerfile:1

FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o sniffer .

FROM alpine:3.20
WORKDIR /
COPY --from=build /app/sniffer /sniffer
RUN mkdir /logs
VOLUME ["/logs"]
EXPOSE 80
ENTRYPOINT ["/sniffer"]
