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

COPY --from=builder /go/bin/mapshare /usr/local/bin/mapshare

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/mapshare"]
CMD ["--"]

# vim: set ft=Dockerfile:
