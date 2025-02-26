FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY . .

RUN go build wiredb.go


FROM alpine:latest

LABEL maintainer="ding@ibyte.me"

WORKDIR /tmp/wiredb

COPY --from=builder /app/wiredb /usr/local/bin/wiredb

COPY config.yaml /tmp/wiredb/config.yaml

EXPOSE 2668

CMD ["/usr/local/bin/wiredb", "--config", "/tmp/wiredb/config.yaml"]