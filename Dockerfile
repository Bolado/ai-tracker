# Use the official Golang image as the base image
FROM golang:1.22 as builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules manifests and download the modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Install templ and generate templates
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN templ generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ai-tracker

# Use a minimal base image for the final image
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built Go application and necessary files from the builder stage
COPY --from=builder /app/ai-tracker /app/websites.json /app/words.json ./
COPY --from=builder /app/website/static ./website/static

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./ai-tracker"]
