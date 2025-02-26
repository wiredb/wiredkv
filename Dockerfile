FROM golang:1.20-alpine AS builder

LABEL maintainer="ding@ibyte.me"

WORKDIR /app

COPY . .
RUN go build wiredb.go

FROM alpine:latest

WORKDIR /data/wiredb

COPY --from=builder /app/wiredb /usr/local/bin/wiredb

COPY config.yaml /data/wiredb/config.yaml

EXPOSE 2668

CMD ["/usr/local/bin/wiredb", "--config", "/data/wiredb/config.yaml"]

# docker build -t wiredb:latest -t wiredb:0.1.1 .