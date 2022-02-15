FROM golang:1.17.7-bullseye AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY CHANGELOG.md /CHANGELOG.md

RUN go mod download

COPY *.go ./
COPY commands/*.go ./commands/
COPY config/*.go ./config/
COPY handlers/*.go ./handlers/
COPY kyDB/*.go ./kyDB/
COPY component/*.go ./component/
COPY update/*.go ./update/

RUN go build -o /kybot

##
## Deploy
##

FROM alpine:latest AS run

WORKDIR /

COPY --from=build /kybot /kybot
COPY --from=build /CHANGELOG.md /CHANGELOG.md
RUN mkdir -p /data

ENTRYPOINT ["/kybot"]