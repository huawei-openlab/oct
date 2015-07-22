# Build Sysbench image
FROM ubuntu:14.04
MAINTAINER Wangkirin <wangqilin2@huawei.com>
ENV DEBIAN_FRONTEND noninteractive


# Install sysbench
RUN sudo apt-get update
RUN sudo apt-get install sysbench -y
VOLUME /data

#Config and Run sysbench test shell
ADD ./fileio.sh /fileio.sh
RUN chmod 755 /fileio.sh
CMD ["/fileio.sh"]

