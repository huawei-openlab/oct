# Build Sysbench image
FROM ubuntu:14.04
MAINTAINER Zhinan Lin (ZenLin) <linzhinan@huawei.com>
ENV DEBIAN_FRONTEND noninteractive


# load  data volume
VOLUME /data

#Add test script
ADD deadloop.sh /deadloop.sh
RUN chmod 755 /deadloop.sh

ENTRYPOINT ["/deadloop.sh"]

