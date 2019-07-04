# Multi-stage build setup (https://docs.docker.com/develop/develop-images/multistage-build/)

# Stage 1 (to create a "build" image, ~850MB)
FROM golang:1.12.6 AS builder
RUN go version

ENV SOURCES /go/src/github.com/nonamenix/sentry2prometheus/

COPY . $SOURCES
WORKDIR $SOURCES

RUN CGO_ENABLED=0 GO111MODULE=on GOOS=linux GOARCH=amd64 go get
RUN go build

RUN ls -la
RUN ls app

# Stage 2 (to create a downsized "container executable")

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder $SOURCES/app .

EXPOSE 9412
ENTRYPOINT ["./sentry2prometheus"]
