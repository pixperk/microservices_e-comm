# Stage 1: Build
FROM golang:1.24-alpine AS build

ENV CGO_ENABLED=1 GOOS=linux GOARCH=amd64

RUN apk --no-cache add gcc g++ make ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
COPY vendor ./vendor
COPY catalog ./catalog

RUN go build -mod=vendor -o /app/catalog ./catalog/cmd/catalog

# Stage 2: Runtime
FROM alpine:3.18

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY --from=build /app/catalog .

EXPOSE 8080

CMD ["./catalog"]
