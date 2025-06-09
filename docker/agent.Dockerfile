FROM golang:1.24.3 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o agent ./cmd/agent/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/agent /app/agent

RUN chmod +x /app/agent

CMD ["sh", "-c", "tail -f /dev/null"]