# Use a lightweight, official Go image as the base
FROM golang:1.22

# Set the working directory inside the container
WORKDIR /app

# Copy go module files for dependency management
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port your application listens on (adjust if necessary)
EXPOSE 8066

# Command to run when the container starts
CMD ["./main"]
