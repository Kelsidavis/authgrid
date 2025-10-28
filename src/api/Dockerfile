# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and source files
COPY go.mod ./
COPY *.go ./

# Download dependencies (this creates go.sum automatically)
RUN go mod download

# Tidy up dependencies
RUN go mod tidy

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o authgrid-api .

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/authgrid-api .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./authgrid-api"]
