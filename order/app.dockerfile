# Stage 1: Build
FROM golang:1.24-alpine AS build

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /app

# Copy dependencies and sources
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY account ./account
COPY catalog ./catalog
COPY order ./order

# Build the order service
RUN go build -mod=vendor -o /app/order ./order/cmd/order

# Stage 2: Runtime
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the compiled binary
COPY --from=build /app/order .

EXPOSE 8080

CMD ["./order"]
