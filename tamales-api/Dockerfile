#FROM golang:1.12.7-alpine
#
#WORKDIR /tamales-api
#COPY . .
#
#RUN go build -o tamales-api *.go

FROM golang:1.12.7 AS build

WORKDIR /tamales-api
COPY . .

RUN go build -o tamales-api *.go

FROM alpine:latest AS runtime
#this seems dumb, but the libc from the build stage is not the same as the alpine libc
#create a symlink to where it expects it since they are compatable. https://stackoverflow.com/a/35613430/3105368
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /tamales-api

EXPOSE 8080

ENV CASSANDRA_HOST cassandra

CMD ["./tamales-api"]
