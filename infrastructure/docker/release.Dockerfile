FROM golang:1.12.7-alpine

ENV GO111MODULE on
EXPOSE 8080

RUN apk update && apk add --virtual build-dependencies build-base
RUN apk add git && apk add ca-certificates

WORKDIR /go/src/bitbucket.org/siimpl/esp-betting/betting-feed

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN make go-bindata
RUN make optimized-docker-service-build

ENTRYPOINT ["/bin/betting-feed-service"]
