#FROM golang:buster AS builder
FROM golang:1.16.6-alpine AS builder

ARG BUILD_VERSION=local
ENV BUILD_VERSION=${BUILD_VERSION}

WORKDIR /app

ADD . .

RUN go build -o /usr/local/bin/undina

FROM alpine:3.14.0

WORKDIR /app

COPY --from=builder /usr/local/bin/undina /app/undina
#COPY --from=builder /app/config.json /app/config.json

# need to fix time zone info in alpine docker image (http://www.csyangchen.com/go-alpine-time-location.html)
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /opt/zoneinfo.zip
ENV ZONEINFO /opt/zoneinfo.zip

EXPOSE 5000

CMD ["/app/undina"]
