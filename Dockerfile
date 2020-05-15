FROM golang:1.12.7-stretch as builder

COPY go.mod /go/src/github.com/pnetwork/sre.ceph.init/go.mod
COPY go.sum /go/src/github.com/pnetwork/sre.ceph.init/go.sum

# Run golang at any directory, not neccessary $GOROOT, $GOPATH
ENV GO111MODULE=on
WORKDIR /go/src/github.com/pnetwork/sre.ceph.init

# RUN go mod init github.com/pnetwork/sre.monitor.metrics
RUN go mod download
COPY cmd /go/src/github.com/pnetwork/sre.ceph.init/cmd
COPY internal /go/src/github.com/pnetwork/sre.ceph.init/internal
#COPY pkg /go/src/github.com/pnetwork/sre.ceph.init/pkg

# Build the Go app
RUN env GOOS=linux GOARCH=amd64 go build -o ceph-init -v -ldflags "-s" github.com/pnetwork/sre.ceph.init/cmd

##### To reduce the final image size, start a new stage with alpine from scratch #####
#FROM alpine:3.9
#RUN apk --no-cache add ca-certificates libc6-compat

# Run as root
#WORKDIR /root/

# Copy the pre-built binary file from the previous stage
#COPY --from=builder /go/src/github.com/pnetwork/sre.monitor.metrics/marvin-exporter /usr/local/bin/marvin-exporter

# EXPOSE 9987

#ENTRYPOINT [ "marvin-exporter" ]
#~




FROM ubuntu:xenial

RUN apt-get update && apt-get install -y \
    apt-transport-https \
    git \
    software-properties-common \
    uuid-runtime \
    wget

ARG CEPH_REPO_URL=https://download.ceph.com/debian-luminous/
RUN wget -q -O- 'https://download.ceph.com/keys/release.asc' | apt-key add -
RUN apt-add-repository "deb ${CEPH_REPO_URL} xenial main"

RUN add-apt-repository -y ppa:gophers/archive

RUN apt-get update && apt-get install -y \
    ceph \
    libcephfs-dev \
    librados-dev \
    librbd-dev \
    golang-1.10-go

# add user account to test permissions
RUN groupadd -g 1010 bob
RUN useradd -u 1010 -g bob -M bob

ENV GOPATH /go
WORKDIR /go/src/github.com/pnetwork/sre.ceph.init
VOLUME /go/src/github.com/ceph/pnetwork/sre.ceph.init

COPY --from=builder /go/src/github.com/pnetwork/sre.ceph.init/ceph-init /usr/local/bin/ceph-init
ENTRYPOINT [ "ceph-init" ]