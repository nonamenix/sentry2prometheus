# Multi-stage build setup (https://docs.docker.com/develop/develop-images/multistage-build/)

# Stage 1 (to create a "build" image, ~850MB)
FROM golang:1.12.6 AS builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

# Stage 2 (to create a downsized "container executable")

FROM alpine:3.7
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/sentry2prometheus /root
COPY --from=builder /app/entrypoint.sh /root
RUN ls -la /root
EXPOSE 9412
ENTRYPOINT ["./entrypoint.sh"]
CMD []
