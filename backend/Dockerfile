# Use the official Golang image to create a build artifact.
# This is known as a multi-stage build.
FROM golang:1.23-bookworm AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# We want to populate the module cache based on the go.mod file to speed up subsequent builds.
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
# -o out/main: specifies the output file name
# CGO_ENABLED=0: disables CGO to create a statically linked binary
# GOOS=linux: specifies the target operating system
RUN CGO_ENABLED=0 GOOS=linux go build -o out/main .

# Start a new stage from scratch for a smaller image
FROM debian:bookworm-slim

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/out/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
