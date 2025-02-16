# Stage 1: Build the Go app
FROM golang:1.24 AS builder

# Set the current working directory inside the container
WORKDIR /chatapp

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN GOOS=$(go env GOOS) GOARCH=$(go env GOARCH) go build -o main .

# Stage 2: Create a minimal image to run the Go app
FROM alpine:latest

# Set the current working directory inside the container
WORKDIR /chatapp

# Copy the binary from the builder stage
COPY --from=builder /chatapp .

# Ensure the executable has the correct permissions (for Unix systems)
RUN chmod +x main

# Expose port 8080 to the outside world
EXPOSE 8080

# Entrypoint for app
ENTRYPOINT [ "./main || ./main.exe" ]
