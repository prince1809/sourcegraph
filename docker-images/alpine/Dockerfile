# This Dockerfile defines the sourcegraph/alpine Docker image, which is the
# base image used by all Sourcegraph Docker images.

FROM alpine:3.9@sha256:644fcb1a676b5165371437feaa922943aaf7afcfa8bfee4472f6860aad1ef2a0

LABEL org.opencontainers.image.url=https://sourcegraph.com/
LABEL org.opencontainers.image.source=https://github.com/prince1809/sourcegraph/
LABEL org.opencontainers.image.documentation=https://docs.sourcegraph.com/

# Add the sourcegraph group, user, and create the home directory.
#
# We use a static GID/UID assignment to ensure files can be chown'd to this
# user on the host machine (where this user does not exist).
RUN addgroup -g 101 -S sourcegraph && adduser -u 100 -S -G sourcegraph -h /home/sourcegraph sourcegraph

# Install bind-tools to ensure working DNS on user-defined Docker networks.
# hadolint ignore=DL3018
RUN apk add --no-cache bind-tools

# Install other packages that are desirable in ALL Sourcegraph Docker images.
# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates mailcap tini
