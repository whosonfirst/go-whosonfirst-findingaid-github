FROM golang:1.16-alpine as builder

RUN mkdir /build

COPY . /build/go-whosonfirst-findingaid-github

RUN apk update && apk upgrade \
    && apk add make libc-dev gcc git \
    && cd /build/go-whosonfirst-findingaid-github \
    && go build -mod vendor -o /usr/local/bin/populate cmd/populate/main.go    

FROM alpine:latest

COPY --from=builder /usr/local/bin/populate /usr/local/bin/

RUN apk update && apk upgrade \
    && apk add ca-certificates