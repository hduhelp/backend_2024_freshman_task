FROM golang:1.23.1 AS builder

WORKDIR /app

COPY . .

RUN go build -o myapp

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/myapp /app/
COPY config.toml .
COPY entrypoint.sh .
RUN chmod +x entrypoint.sh

ENTRYPOINT ["./entrypoint.sh"]

EXPOSE 8080