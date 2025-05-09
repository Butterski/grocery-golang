FROM golang:1.24 AS build
WORKDIR /go/src
COPY go ./go
COPY main.go .
COPY go.mod .
COPY go.sum .

# Install dependencies
RUN go mod download
RUN go mod tidy

ENV CGO_ENABLED=0

# Build with extra error information
RUN go build -o openapi .

# Use a more feature-rich base image to support networking with PostgreSQL
FROM alpine:3.21 AS runtime
# Install CA certificates for secure connections
RUN apk --no-cache add ca-certificates

WORKDIR /app
ENV GIN_MODE=release

# Copy the binary from the build stage
COPY --from=build /go/src/openapi ./

# Copy the API specification directory
COPY api ./api

# Expose the API port
EXPOSE 8080/tcp

# Run the application
ENTRYPOINT ["./openapi"]
