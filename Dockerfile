FROM golang:1.23.2-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod downloand

COPY ./ ./
RUN go build ./cmd/server

FROM alpine:3.20

COPY --from=buidler /app/server /usr/local/bin
CMD ["server"]
