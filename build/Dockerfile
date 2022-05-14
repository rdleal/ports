FROM golang:1.17-alpine AS builder

RUN apk update
RUN apk --no-cache add build-base git openssh-client ca-certificates

WORKDIR /app

COPY go.mod .

RUN go mod download -x

ADD . /ports-service

RUN cd /ports-service/cmd/api && \
    go build -tags musl -buildmode=pie -o ports-service-api

RUN cp /ports-service/cmd/api/ports-service-api  /app/ports-service-api

FROM alpine:latest

RUN apk --no-cache add tzdata

WORKDIR /app

COPY --from=builder /app/ports-service-api .

ENV PORT 8080

ENTRYPOINT /app/ports-service-api
