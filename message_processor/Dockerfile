FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY . ./
COPY go.mod ./
COPY go.sum ./


RUN go build -o message_processor

FROM alpine:3.22
WORKDIR /app
COPY .env ./
COPY latest_timestamp.json ./

# Install system dependencies: ca-certificates, yt-dlp, python3, py3-pip, and build dependencies for pip/uv
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/message_processor .

ENTRYPOINT ["./message_processor"]