# Use the official Golang image as a build stage
FROM golang:1.23-alpine AS builder

# Set the working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the entire project
COPY . .

# Build the Go application
RUN go build -o hospital-system .

# Use a minimal base image for running the application
FROM alpine:latest

# Set environment variables
ENV PORT=8080

# Install CA certificates for HTTPS
RUN apk add --no-cache ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/hospital-system .

# Expose the application's port
EXPOSE 8080

# Run the application
CMD ["./hospital-system"]
