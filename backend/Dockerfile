# Use an official Golang image as a build stage
FROM golang:1.21-alpine AS builder

# Set the Go module proxy
ENV GOPROXY=https://goproxy.io,direct

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy the application code
COPY backend/ .

# Build the binary
RUN go build -o main .

# Final stage for a minimal image
FROM alpine:3.18

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Expose the port
EXPOSE 1323

# Run the application
CMD ["./main"]
