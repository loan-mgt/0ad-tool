# --- Build stage ---
FROM golang:1.24.3-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY api/go.mod api/go.sum ./
RUN go mod download

# Copy the source code
COPY api/. .

# Build the Go app
RUN go build -o app .

# --- Final stage ---
FROM alpine:latest

WORKDIR /app

# Copy the built binary
COPY --from=builder /app/app .

# Copy the templates directory
COPY  templates /templates
COPY  public /app/public

# Expose port (change if needed)
EXPOSE 8081

# Run the app
CMD ["./app"]
