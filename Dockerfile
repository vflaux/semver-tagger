# Copyright 2022 vflaux

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#  		http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARG ALPINE_VERSION=3.12

FROM gcr.io/distroless/static

COPY bin/ /usr/local/bin

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
