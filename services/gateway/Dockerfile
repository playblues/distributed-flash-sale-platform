FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o gateway-app ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gateway-app .
EXPOSE 3000
CMD ["./gateway-app"]
