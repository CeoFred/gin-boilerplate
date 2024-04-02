# Use a smaller base image (Alpine) for reduced size
FROM golang:alpine as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app with static linking for reduced dependencies
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Start a new stage for the final minimal image
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/app .

# Expose port 3006 to the outside world
EXPOSE 3006

# Command to run the executable
CMD ["./app"]
