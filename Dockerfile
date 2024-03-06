## build
FROM golang:1.22.1-alpine3.19 AS build-env

RUN apk add build-base

ADD . /go/src/app

WORKDIR /go/src/app

RUN go env -w GO111MODULE=on \
    && go mod tidy \
    && cd ./cmd/server \
    && go build -o ../../relay

## run
FROM alpine:3.19

RUN mkdir -p /ep && mkdir -p /ep/log

WORKDIR /ep

COPY --from=build-env /go/src/app /ep/

ENV PATH $PATH:/aa

EXPOSE 80
CMD ["/ep/relay"]
