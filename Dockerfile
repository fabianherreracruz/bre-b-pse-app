# Build stage
FROM golang:1.21 AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source
COPY . .

# Build app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bre-b-pse-app ./cmd

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/bre-b-pse-app .

# Expose port
EXPOSE 8080

# Run app
CMD ["./bre-b-pse-app"]
