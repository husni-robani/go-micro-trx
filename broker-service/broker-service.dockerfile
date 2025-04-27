# Stage 1
FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN env CGO_ENABLED=0 GOOS=linux go build -o brokerApp ./cmd/api

# Stage 2
FROM alpine:latest

WORKDIR /app

COPY --from=build /app/brokerApp .

CMD [ "/app/brokerApp" ]
