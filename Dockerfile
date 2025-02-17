# Stage 1: Build the Go app
FROM golang:1.24-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum before the rest to take advantage of Docker caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the Go app
RUN go build -o main ./cmd

# Stage 2: Create a minimal image to run the Go app
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Install necessary dependencies
RUN apk add --no-cache ca-certificates

# Copy the built Go app from the builder stage
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/public .

# Command to run the Go app
CMD ["./main"]
