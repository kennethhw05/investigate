FROM golang:1.12.4

WORKDIR /betting-feed
RUN go get github.com/derekparker/delve/cmd/dlv

RUN ["go", "get", "github.com/githubnemo/CompileDaemon"]
