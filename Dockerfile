# Start from a Go Alpine base image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go code to the container
COPY . .
RUN go get main
# Build the Go application
RUN go build -o main .

# Set the command to run when the container starts
CMD ["./main"]
