# Stage 1: Build
FROM golang:1.24-alpine AS build

# Set necessary Go env vars
ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# Install build tools
RUN apk --no-cache add gcc g++ make ca-certificates

# Set workdir
WORKDIR /app

# Copy dependencies
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY account ./account
COPY catalog ./catalog
COPY order ./order
COPY graphql ./graphql

# Build the binary
RUN go build -mod vendor -o app ./graphql

# Stage 2: Run
FROM alpine:3.18

# Add CA certs (just in case)
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy binary
COPY --from=build /app/app .

# Expose port
EXPOSE 8080

# Run the app
CMD ["./app"]
