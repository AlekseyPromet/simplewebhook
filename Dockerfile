FROM golang:1.22.0-alpine3.19 AS builder

ARG ALPINE_VER

RUN apk add tzdata
RUN apk add --no-cache ca-certificates git

WORKDIR /simplewebhook
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum

# install updates and build executable
RUN apk update -X http://dl-3.alpinelinux.org/alpine/v3.19/main && \
    apk upgrade --no-cache -X http://dl-3.alpinelinux.org/alpine/v3.19/main && \
    CGO_ENABLED=0 GOOS=linux go build -mod=readonly -a -installsuffix cgo -o app  \
    -ldflags "-X 'main.xBuildVersion=${VERSION}' -X 'main.xBuildHashCommit=$HASHCOMMIT'"  ./cmd

# copy to alpine image
FROM alpine:3.19

# create user other than root and install updated
RUN addgroup -g 101 app && \
    adduser -H -u 101 -G app -s /bin/sh -D app && \
    apk update --no-cache -X http://dl-3.alpinelinux.org/alpine/v3.19/main && \
    apk upgrade --no-cache -X http://dl-3.alpinelinux.org/alpine/v3.19/main/main

# place all necessary executables and other files into /app directory
WORKDIR /app/
COPY --from=builder --chown=app:app /src/app .

# run container as new non-root user
USER app

CMD ["/app/app"]