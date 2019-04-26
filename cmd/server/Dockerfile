FROM alpine:3.9@sha256:644fcb1a676b5165371437feaa922943aaf7afcfa8bfee4472f6860aad1ef2a0 AS ctags

RUN apk --no-cache add --virtual build-deps curl jansson-dev \
  libseccomp-dev linux-headers autoconf pkgconfig make automake \
  gcc g++ binutils