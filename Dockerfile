# Use a lightweight Go image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod .

RUN go mod download || true

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o go-web-service ./main.go

# Expose port 8080 (same as in main.go)
EXPOSE 8080

# Run the compiled Go application
CMD ["/app/go-web-service"]