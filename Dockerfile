# syntax=docker.io/docker/dockerfile:1.4

FROM scratch as base
ARG TARGETOS
ARG TARGETARCH
COPY bin/semver-tagger_${TARGETOS}_${TARGETARCH} /usr/local/bin/semver-tagger
ENTRYPOINT ["/usr/local/bin/semver-tagger"]

FROM debian:bookworm-slim as debian
ARG TARGETOS
ARG TARGETARCH
COPY --link bin/semver-tagger_${TARGETOS}_${TARGETARCH} /usr/local/bin/semver-tagger

FROM debian as github-action
COPY --link semver-tagger-action/entrypoint.sh /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

FROM alpine:3.17 as alpine
ARG TARGETOS
ARG TARGETARCH
COPY --link bin/semver-tagger_${TARGETOS}_${TARGETARCH} /usr/local/bin/semver-tagger