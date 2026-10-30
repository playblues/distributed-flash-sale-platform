FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o inventory-app ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/inventory-app .
EXPOSE 3001 50051
CMD ["./inventory-app"]
