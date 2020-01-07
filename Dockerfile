FROM alpine as prefetch

COPY . /src
WORKDIR /src

RUN set -ex \
 && apk add \
      bash \
      curl \
      make \
 && make assets


FROM golang:alpine as builder

COPY . /go/src/github.com/Luzifer/mapshare
WORKDIR /go/src/github.com/Luzifer/mapshare

RUN set -ex \
 && apk add --update git \
 && go install \
      -ldflags "-X main.version=$(git describe --tags --always || echo dev)" \
      -mod=readonly

FROM alpine:latest

LABEL maintainer "Knut Ahlers <knut@ahlers.me>"

RUN set -ex \
 && apk --no-cache add \
      ca-certificates

COPY --from=builder   /go/bin/mapshare  /usr/local/bin/mapshare
COPY --from=prefetch  /src/frontend     /usr/local/share/mapshare/frontend

EXPOSE 3000
WORKDIR /usr/local/share/mapshare

ENTRYPOINT ["/usr/local/bin/mapshare"]
CMD ["--"]

# vim: set ft=Dockerfile:
