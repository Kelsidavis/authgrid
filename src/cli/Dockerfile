# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and source files
COPY go.mod ./
COPY main.go ./

# Download dependencies
RUN go mod download

# Build the CLI
RUN CGO_ENABLED=0 GOOS=linux go build -o authgrid-cli .

# Create a minimal runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/authgrid-cli .

# Set as entrypoint
ENTRYPOINT ["/app/authgrid-cli"]
CMD ["help"]
