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

VENDOR ?= vflaux
SOURCE ?= https://github.com/vflaux/semver-tagger
IMAGE_NAME := semver-tagger
REMOTE_NAME := $(REGISTRY_URL)/$(IMAGE_NAME)
LOCAL_NAME := build/$(IMAGE_NAME)
IMAGE_VERSION ?= test
IMAGE_VERSION := $(IMAGE_VERSION:v%=%)

export CGO_ENABLED=0
export DOCKER_BUILDKIT=1

default: build

build:
	go build -o bin/semver-tagger main.go

test: fmt vet
	go test ./... -coverprofile cover.out

run: fmt vet
	go run ./main.go

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -f ../bin/semver-tagger

image: build
	docker build . \
	  --build-arg VENDOR="${VENDOR}" \
	  --build-arg SOURCE="${SOURCE}" \
		--build-arg REVISION="${shell git rev-parse HEAD}" \
		--build-arg VERSION="${IMAGE_VERSION}" \
		--build-arg TITLE="${IMAGE_NAME}" \
		--build-arg CREATED="$(date +%Y-%m-%dT%H:%M:%S%z)" \
	  -t "${LOCAL_NAME}:${IMAGE_VERSION}"
