FROM golang:1.19beta1-alpine as builder

RUN apk add chromium

WORKDIR /usr/src/ihysApp

COPY . .

RUN go get github.com/mafredri/cdp

RUN go build -o bin/ihysApp cmd/main/main.go

CMD ["./bin/ihysApp"]
