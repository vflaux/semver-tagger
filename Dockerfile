ARG ALPINE_VERSION=3.17

FROM alpine:${ALPINE_VERSION}

RUN apk add make --no-cache

ARG TARGETOS
ARG TARGETARCH

COPY --chmod=555 bin/semver-tagger_${TARGETOS}_${TARGETARCH} /usr/local/bin/semver-tagger

ARG VENDOR
ARG SOURCE
ARG VERSION
ARG REVISION
ARG TITLE
ARG CREATED

LABEL \
    org.opencontainers.image.description="A tool to remotely tag an image if it is the latest version" \
    org.opencontainers.image.vendor="$VENDOR" \
    org.opencontainers.image.source="$SOURCE" \
    org.opencontainers.image.version="$VERSION" \
    org.opencontainers.image.revision="$REVISION" \
    org.opencontainers.image.title="$TITLE" \
    org.opencontainers.image.created="$CREATED"
