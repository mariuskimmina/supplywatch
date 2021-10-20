# Start from golang base image
FROM golang:1.13-alpine as builder

# Set the current working directory inside the container
WORKDIR /build

# Copy go.mod, go.sum files and download deps
COPY go.mod go.sum ./
RUN go mod download

# Copy sources to the working directory
COPY . .

# Build th Go app
ARG project
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -a -v -o service $project

# Start a new stage from busybox
FROM busybox:latest

WORKDIR /dist

# Copy the build artifacts from the previous stage
COPY --from=builder /build/service .

# Run the executable
CMD ["./service"]
