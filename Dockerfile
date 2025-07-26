FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -o sniffer ./cmd/main.go

FROM gcr.io/distroless/static-debian11 AS runner

WORKDIR /app

COPY --from=builder /app/sniffer /app/sniffer

EXPOSE 3185

CMD ["/app/tweakio"]