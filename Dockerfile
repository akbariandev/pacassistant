# Modify the Dockerfile to copy uploaded files during build
# Use the official Golang image as the base image
FROM golang:1.21.4 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install Go module dependencies
RUN go mod download

# Copy the rest of the application code to the working directory
COPY . .

# Build the Golang application
RUN go build -o main .

# Start a new stage to create a smaller image
FROM golang:1.21.4

# Set the working directory inside the container
WORKDIR /app

# Copy the binary from the builder stage to the new stage
COPY --from=builder /app/main .

# Expose the port the application will run on
EXPOSE 8050

# Start the Golang application
CMD ["./main","-c", "/etc/bot/config.yml", "run"]