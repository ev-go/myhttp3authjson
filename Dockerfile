FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /go/src/app

COPY . .

RUN go build -o main .

EXPOSE 3000

CMD ["/go/src/app/main"]