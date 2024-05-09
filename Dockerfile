FROM golang:1.21-alpine as builder

RUN apk add --no-cache bash make git curl

RUN mkdir /pacassistant
WORKDIR /pacassistant

COPY . /pacassistant

# build
WORKDIR /pacassistant
RUN make build-bot

FROM alpine:3.14
RUN apk --no-cache add ca-certificates tzdata git
#RUN mkdir -p /etc/bot
RUN mkdir /pacassistant
COPY --from=builder /pacassistant/main /pacassistant
COPY --from=builder /pacassistant/config/config.yml /pacassistant/config/config.yml

RUN chmod +x /pacassistant/main
CMD ["./pacassistant/main", "-c", "./pacassistant/config/config.yml", "run"]
