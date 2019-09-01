# Copyright (c) 2019 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

# ------------------------------------------------------------------------------
# Build web assets
# ------------------------------------------------------------------------------
FROM node:10.15.3 as base

ADD web/ /web
WORKDIR /web
ENV CYPRESS_INSTALL_BINARY=0

RUN npm ci --prefer-offline && npm run-script build

# ------------------------------------------------------------------------------
# Install go tools and build binary
# ------------------------------------------------------------------------------
FROM golang:1.12 as builder

WORKDIR /workspace
ADD . /workspace
COPY --from=base /web ./web
ENV GOFLAGS=-mod=vendor GO111MODULE=on

RUN make go-install
RUN go generate ./web
RUN make generate
RUN make octant-dev

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

