FROM golang:1.21-alpine as builder

RUN mkdir /app
WORKDIR /app

COPY . /app

# build
WORKDIR /app
RUN make build-bot

FROM alpine:3.14
RUN apk --no-cache add ca-certificates tzdata git
RUN mkdir -p /etc/bot
RUN mkdir /app
COPY --from=builder /app/main /app
RUN chmod +x /app/main
CMD ["./app/main", "-c", "/etc/bot/config.yml", "run"]
