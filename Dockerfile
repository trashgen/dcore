FROM golang:alpine as builder

ADD / $GOPATH/src/dcore
WORKDIR $GOPATH/src/dcore/cmd/config
RUN go build -o /go/bin/dcore/config .

RUN apk add --no-cache git
RUN go get -d -v github.com/lib/pq
RUN apk del git

WORKDIR $GOPATH/src/dcore/cmd/point
RUN go build -o /go/bin/dcore/point .
WORKDIR $GOPATH/src/dcore/cmd/node
RUN go build -o /go/bin/dcore/node .


FROM alpine
COPY --from=builder /go/bin/dcore/node /go/bin/dcore/node
COPY --from=builder /go/bin/dcore/point /go/bin/dcore/point
COPY --from=builder /go/bin/dcore/config /go/bin/dcore/config