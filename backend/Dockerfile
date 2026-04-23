FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/api ./cmd/api

FROM alpine:3.21

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata && adduser -D -g '' appuser

COPY --from=builder /bin/api /app/api
COPY migration /app/migration

USER appuser

EXPOSE 8080

CMD ["/app/api"]
