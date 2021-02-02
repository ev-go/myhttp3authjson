FROM golang:alpine

WORKDIR /mnt/app
COPY . .

EXPOSE 3000