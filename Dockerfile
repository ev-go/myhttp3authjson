FROM golang:alpine

WORKDIR /go/src/app
COPY . .

RUN /bin/sh -c 'go run main.go; \
root; \
1'


EXPOSE 3000


