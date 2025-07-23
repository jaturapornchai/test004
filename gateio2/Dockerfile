# Use official golang runtime as base image
FROM golang:1.24-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Install necessary packages
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Copy .env file 
COPY .env .env

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateio-trading-bot .

# Start a new stage from scratch
FROM alpine:latest

# Install ca-certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata

# Set timezone to Bangkok
ENV TZ=Asia/Bangkok

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/gateio-trading-bot .

# Copy .env file
COPY --from=builder /app/.env .

# Copy prompts directory
COPY --from=builder /app/prompts ./prompts/

# Create logs directory
RUN mkdir -p logs

# Make sure the binary is executable
RUN chmod +x ./gateio-trading-bot

# Expose port (if needed for health checks)
EXPOSE 8080

# Command to run the executable  
CMD ["./gateio-trading-bot"]
