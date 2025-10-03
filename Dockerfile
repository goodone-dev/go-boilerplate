# Use the official Golang image as a base image
FROM golang:1.24-alpine

# Set the current working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main ./cmd/api/main.go

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]
