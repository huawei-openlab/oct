# Build Sysbench image
FROM ubuntu:14.04
MAINTAINER Wangkirin <wangqilin2@huawei.com>
ENV DEBIAN_FRONTEND noninteractive


# Install sysbench
RUN sudo apt-get update
RUN sudo apt-get install sysbench -y
VOLUME /data

#Config and Run sysbench test shell
ADD memoryrw.sh /memoryrw.sh
RUN chmod 755 /memoryrw.sh
CMD ["/memoryrw.sh"]

