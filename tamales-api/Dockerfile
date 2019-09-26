FROM golang:1.12.7 AS build

WORKDIR /tamales-api
COPY . .
RUN CGO_ENABLED=0 go build -o tamales-api *.go

FROM alpine:latest AS runtime

WORKDIR /
COPY --from=0 /tamales-api/tamales-api /tamales-api
EXPOSE 8080
ENV CASSANDRA_HOST cassandra
CMD ["/tamales-api"]
