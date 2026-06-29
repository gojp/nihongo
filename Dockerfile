# Stage 1: Build the binary
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Stage 2: Final minimal image
FROM scratch
WORKDIR /root/
# Copy the binary from the builder stage
COPY --from=builder /app/main .
# Expose the port your app listens on
EXPOSE 8080
# Run the binary
CMD ["./main", "-addr=0.0.0.0:8080"]
