FROM golang:1.19 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /build/
COPY . /build/

RUN go build -tags="timetzdata" -ldflags "-w -s" ./cmd/benthos

FROM jeffail/benthos

LABEL maintainer="Antoine Girard <antoine.girard@sapk.fr>"

# replace original benthos
COPY --from=build /build/benthos .

