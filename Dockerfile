# Start from a Go Alpine base image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go code to the container
COPY main.go .

# Build the Go application
RUN go build -o main .

# Expose the port the application listens on
EXPOSE 8080

# Set the command to run when the container starts
CMD ["./main"]
