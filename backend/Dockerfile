# Base image with Go 1.23.11 and Alpine
FROM golang:1.23.11-alpine

# Install Python and tools
RUN apk add --no-cache \
    python3 \
    bash \
    build-base 

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Set up Go modules
RUN go mod tidy

ENV CGO_ENABLED=1
RUN go build -o main .
RUN ls -lh .

# Run TLS/env initialization
RUN python3 init.py

# Expose backend port
EXPOSE 8080

# Start backend server
CMD ["./main"]
