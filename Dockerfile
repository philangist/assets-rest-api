FROM golang:1.10

RUN mkdir -p /go/src/github.com/philangist/frameio-assets
WORKDIR /go/src/github.com/philangist/frameio-assets

ADD . /go/src/github.com/philangist/frameio-assets

RUN go get github.com/lib/pq
