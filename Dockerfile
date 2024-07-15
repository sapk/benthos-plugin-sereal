FROM golang:1.22 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build/
COPY . /build/

RUN go build -tags="timetzdata" -ldflags "-w -s" ./cmd/benthos

FROM ghcr.io/redpanda-data/connect

LABEL maintainer="Antoine Girard <antoine.girard@sapk.fr>"
LABEL org.opencontainers.image.source="https://github.com/sapk/benthos-plugin-sereal"

# replace original benthos binary and configuration
COPY ./config/sereal.yaml /benthos.yaml
COPY --from=build /build/benthos .

