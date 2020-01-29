# Copyright (c) 2019 the Octant contributors. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# ------------------------------------------------------------------------------
# Build web assets
# ------------------------------------------------------------------------------
ARG NODE_VERSION

FROM node:$NODE_VERSION as base

ADD web/ /web
WORKDIR /web
ENV CYPRESS_INSTALL_BINARY=0

RUN npm ci --prefer-offline && npm run-script build

# ------------------------------------------------------------------------------
# Install go tools and build binary
# ------------------------------------------------------------------------------
FROM golang:1.13 as builder

WORKDIR /workspace
ADD . /workspace
COPY --from=base /web ./web
ENV GOFLAGS=-mod=vendor GO111MODULE=on

RUN go run build.go go-install
RUN go generate ./pkg/icon
RUN go generate ./web
RUN go run build.go build

# ------------------------------------------------------------------------------
# Running container
# ------------------------------------------------------------------------------
FROM ubuntu:bionic

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        apt-transport-https \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /workspace/build/octant /octant
RUN chmod +x /octant

RUN useradd -s /sbin/nologin -M -u 10000 -U user
USER user

VOLUME [ "/kube"]
EXPOSE 7777

