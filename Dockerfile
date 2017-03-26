FROM golang:1.8.0-alpine

ADD . /go/src/github.com/Tinker-Ware/digital-ocean-service

RUN go install github.com/Tinker-Ware/digital-ocean-service   

ENTRYPOINT /go/bin/digital-ocean-service

EXPOSE 3000