# Build stage
FROM golang:1.25.3-alpine3.22 AS builder

# Set working directory
WORKDIR /build

# Copy dependencies file
COPY go.mod go.sum* ./

# Install dependencies
RUN go mod download

# Copy runner script and default execute module
COPY main.go execute.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o runner .

# Final stage
FROM alpine:3.22

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy the binary from builder
COPY --from=builder /build/runner /app/runner
RUN chmod +x /app/runner

# Set working directory
WORKDIR /app

# Set entrypoint
ENTRYPOINT ["/app/runner"]