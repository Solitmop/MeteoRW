# Use multi-stage build to reduce image size
FROM golang:1.21-alpine AS builder

# Install git (needed for go modules)
RUN apk add --no-cache git

# Set the working directory inside the container
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o meteodata-service .

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/meteodata-service .

# Expose port
EXPOSE 8080

# Run the application
CMD ["./meteodata-service"]