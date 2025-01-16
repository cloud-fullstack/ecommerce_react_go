# Use an official Golang image as a build stage
FROM golang:1.21-alpine AS builder

# Set the Go module proxy
ENV GOPROXY=https://goproxy.io,direct

# Set the working directory
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies and tidy the module
RUN go mod tidy

# Copy the application code
COPY . .

# Build the binary
RUN go build -o main .

# Final stage for a minimal image
FROM alpine:3.18

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the config.sh file (if used)
COPY config.sh /app/config.sh

# Make the config.sh script executable
RUN chmod +x /app/config.sh

# Define build-time arguments
ARG API_PORT
ARG STATIC_PATH
ARG REACT_APP_DOMAIN_NAME
ARG REACT_APP_API_URL
ARG DB_HOST
ARG DB_PORT
ARG DB_USER
ARG DB_NAME
ARG DB_SSLMODE

# Set environment variables
ENV API_PORT=${API_PORT}
ENV STATIC_PATH=${STATIC_PATH}
ENV REACT_APP_DOMAIN_NAME=${REACT_APP_DOMAIN_NAME}
ENV REACT_APP_API_URL=${REACT_APP_API_URL}
ENV DB_HOST=${DB_HOST}
ENV DB_PORT=${DB_PORT}
ENV DB_USER=${DB_USER}
ENV DB_NAME=${DB_NAME}
ENV DB_SSLMODE=${DB_SSLMODE}

# Expose the port
EXPOSE 1323

# Run the application
CMD ["./main"]