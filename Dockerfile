FROM golang:1.12-alpine as build-env

RUN mkdir /simple-golang-crawler
WORKDIR /simple-golang-crawler
COPY go.mod .
COPY go.sum .
ENV GOPROXY="https://goproxy.io" GO111MODULE=on
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/simple-golang-crawler cmd/start-concurrent-engine.go
FROM alpine:3.7
RUN apk update \
        && apk upgrade \
        && apk add --no-cache \
        ca-certificates \
        && update-ca-certificates 2>/dev/null || true \
        && apk add --no-cache ffmpeg
COPY --from=build-env /go/bin/simple-golang-crawler /go/bin/simple-golang-crawler
ENTRYPOINT ["/go/bin/simple-golang-crawler"]

