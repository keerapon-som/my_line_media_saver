FROM golang:1.24.3-alpine AS builder
WORKDIR /app
COPY . ./
COPY go.mod ./
COPY go.sum ./


RUN go build -o line_webhook

FROM alpine:3.22
WORKDIR /app
COPY .env ./

# Install system dependencies: ca-certificates, yt-dlp, python3, py3-pip, and build dependencies for pip/uv
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/line_webhook .

ENTRYPOINT ["./line_webhook"]