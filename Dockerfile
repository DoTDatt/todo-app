FROM golang:1.26.2-alpine3.23 AS builder

WORKDIR /Build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:3.23 
WORKDIR /App

COPY --from=builder /Build/main .
ENV APP_ENV=production

EXPOSE 8080

CMD ["./main"]