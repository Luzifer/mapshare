FROM golang:alpine as builder

COPY . /go/src/github.com/Luzifer/mapshare/frontend
WORKDIR /go/src/github.com/Luzifer/mapshare/frontend

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

COPY --from=builder /go/bin/frontend /usr/local/bin/frontend

EXPOSE 3000

ENTRYPOINT ["/usr/local/bin/frontend"]
CMD ["--"]

# vim: set ft=Dockerfile:
