# syntax=docker/dockerfile:1

FROM golang:1.20 as build

ENV CGO_ENABLED 0
ENV GO_API_PORT 5000

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app

FROM alpine:3.16.0 as production
# Add certificates
RUN apk add --no-cache ca-certificates
# Copy built binary from builder
COPY --from=build app .
# Expose port
EXPOSE ${GO_API_PORT}
# Exec built binary
CMD ./app