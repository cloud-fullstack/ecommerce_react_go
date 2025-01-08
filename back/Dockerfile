# Use an official Golang image as a build stage
FROM golang:1.21-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy the Go module files and install dependencies
COPY back/go.mod back/go.sum ./

# Download dependencies (no need to specify them individually)
RUN go mod download

# Copy the application code (from the 'back' folder)
COPY back/ .

# Build the binary using the vendor directory
RUN go build -o main .

# Final stage for a minimal image
FROM alpine:3.18

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the config.sh file from the 'back' directory
COPY back/config.sh config.sh

# Make the config.sh script executable
RUN chmod +x config.sh

# Set build-time arguments for optional variables
ARG VITE_API_URL
ARG API_URL
ARG DOMAIN_NAME

# Set the environment variables
ENV VITE_API_URL=${VITE_API_URL}
ENV API_URL=${API_URL}
ENV DOMAIN_NAME=${DOMAIN_NAME}

# Expose the port the application will use
EXPOSE 1323

# Run the command to start the application with environment variables loaded
CMD ["sh", "-c", ". ./config.sh && ./main"]
