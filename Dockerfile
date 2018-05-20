FROM golang:1.10

RUN mkdir -p /go/src/github.com/philangist/assets-rest-api
WORKDIR /go/src/github.com/philangist/assets-rest-api

ADD . /go/src/github.com/philangist/assets-rest-api

RUN go get github.com/lib/pq
