FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN env CGO_ENABLED=0 GOOS=linux go build -o transactionApp ./cmd/api

# stage 2
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/transactionApp .

CMD ["/app/transactionApp"]