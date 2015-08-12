FROM ubuntu:14.04

MAINTAINER Tiesheng Cheng <chengtiesheng@huawei.com>

ENV VERSION=3.0.11 \
    BUILD_DEPS="gcc libc6-dev make curl"

RUN apt-get update && \
    apt-get -y install --no-install-recommends $BUILD_DEPS && \
    curl -sL http://downloads.es.net/pub/iperf/iperf-${VERSION}.tar.gz | tar xzf - -C /tmp --strip 1 && \
    (cd /tmp; ./configure) && \
    make -C /tmp && \
    make -C /tmp install && \
    apt-get --purge autoremove -y $BUILD_DEPS && \
    apt-get clean && \
    rm -rf /var/cache/* /var/lib/apt/lists/* /tmp/*

ENTRYPOINT ["/usr/local/bin/iperf3"]
