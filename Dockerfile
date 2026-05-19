# Cloudflare Containers requires linux/amd64
FROM --platform=linux/amd64 golang:1.23-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /server .

FROM --platform=linux/amd64 alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

COPY --from=builder /server .

ENV PORT=8080
EXPOSE 8080

CMD ["./server"]
