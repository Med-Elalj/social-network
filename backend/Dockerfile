# Base image with Go 1.23.11 and Alpine
FROM golang:1.23.11-alpine

# Install Python and tools
RUN apk add --no-cache \
    python3 \
    py3-pip \
    py3-virtualenv \
    git \
    bash \
    build-base \
    curl \
    ca-certificates

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Set up Go modules
RUN go mod tidy

# Set up Python virtual environment and install deps
RUN python3 -m venv /app/venv
ENV PATH="/app/venv/bin:$PATH"
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

ENV CGO_ENABLED=1
RUN go build -o main .
RUN ls -lh .

# Run TLS/env initialization
RUN python3 init.py

# Expose backend port (assuming HTTPS on 8080)
EXPOSE 8080

# Start backend server
CMD ["./main"]
