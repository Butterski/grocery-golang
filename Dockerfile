FROM golang:1.24 AS build
WORKDIR /go/src
COPY go ./go
COPY main.go .
COPY go.sum .
COPY go.mod .

# Install dependencies
RUN go mod download
RUN go mod tidy

ENV CGO_ENABLED=0

RUN go build -o openapi .

# Use a more feature-rich base image to support networking with PostgreSQL
FROM alpine:3.21 AS runtime
# Install CA certificates for secure connections
RUN apk --no-cache add ca-certificates

WORKDIR /app
ENV GIN_MODE=release

# Copy the binary from the build stage
COPY --from=build /go/src/openapi ./

# Expose the API port
EXPOSE 8080/tcp

# Run the application
ENTRYPOINT ["./openapi"]
