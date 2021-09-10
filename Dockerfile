ARG BUILD_VERSION=development
ARG ENV_PHASE=development

FROM golang:1.16.6-alpine AS builder

ARG ENV_PHASE
ENV ENV_PHASE=${ENV_PHASE}

WORKDIR /app

COPY . .

# Override config
RUN cp ./config/${ENV_PHASE}/config.json ./config/config.json

RUN go build -o ./undina

FROM alpine:3.14.0

ARG BUILD_VERSION
ENV BUILD_VERSION=${BUILD_VERSION}

WORKDIR /app

COPY --from=builder /app/undina /app/undina
COPY --from=builder /app/config/config.json /app/config/config.json

# need to fix time zone info in alpine docker image (http://www.csyangchen.com/go-alpine-time-location.html)
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /opt/zoneinfo.zip
ENV ZONEINFO /opt/zoneinfo.zip

EXPOSE 5000

CMD ["/app/undina"]
