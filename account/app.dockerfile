# Stage 1: Build
FROM golang:1.24-alpine AS build

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

# Install build tools
RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /app

# Copy go mod files and vendor
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY account ./account

# Build the account service binary
RUN go build -mod=vendor -o /app/account ./account/cmd/account

# Stage 2: Runtime
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy compiled binary from build stage
COPY --from=build /app/account .

EXPOSE 8080

CMD ["./account"]
