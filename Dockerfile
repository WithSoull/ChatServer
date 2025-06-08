FROM golang:1.23.4-alpine3.21 AS builder

COPY . /chat-server/
WORKDIR /chat-server/

RUN go mod download
RUN go build -o ./bin/chat-server cmd/server/main.go

FROM alpine:3.21
WORKDIR /root/
COPY --from=builder /chat-server/bin/chat-server .
COPY .env .

CMD ["./chat-server"]
