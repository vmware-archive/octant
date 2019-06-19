# Copyright (c) 2019 VMware, Inc. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

FROM golang:1.12 as builder

ADD . /go/src/github.com/vmware/octant
WORKDIR /go/src/github.com/vmware/octant
RUN hacks/setup-docker.sh
RUN make octant-dev

FROM ubuntu:bionic

RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        apt-transport-https \
        ca-certificates \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /go/src/github.com/vmware/octant/build/octant /octant
RUN chmod +x /octant

RUN useradd -s /sbin/nologin -M -u 10000 -U user
USER user

VOLUME [ "/kube"]
EXPOSE 7777
