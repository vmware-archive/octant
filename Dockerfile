FROM golang:1.12 as builder

ADD . /go/src/github.com/heptio/developer-dash
WORKDIR /go/src/github.com/heptio/developer-dash
RUN hacks/setup-docker.sh
RUN make clustereye-dev

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/src/github.com/heptio/developer-dash/build/clustereye /clustereye
RUN chmod +x /clustereye && \
    mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

VOLUME [ "/kube"]
EXPOSE 7777