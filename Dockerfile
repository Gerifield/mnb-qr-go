FROM golang:1.14-alpine AS builder

RUN apk update && apk add ca-certificates

ADD . /go/src/github.com/gerifield/mnb-qr-go
WORKDIR /go/src/github.com/gerifield/mnb-qr-go

ENV CGO_ENABLED=0

RUN go test -v ./...
RUN go build src/cmd/qr-server/qr-server.go


FROM alpine:3.11

COPY --from=builder /go/src/github.com/gerifield/mnb-qr-go/qr-server /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/qr-server"]

EXPOSE 8080