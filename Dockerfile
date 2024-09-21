FROM golang:1.23.1-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server .

from alpine:latest

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080
EXPOSE 1935

CMD ["./server"]
