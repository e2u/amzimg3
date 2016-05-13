FROM ubuntu:latest

MAINTAINER weidewang weidewang@funguide.com.cn

RUN \
  sed -i 's/# \(.*multiverse$\)/\1/g' /etc/apt/sources.list && \
  apt-get update && \
  locale-gen en_US.UTF-8

ENV TZ 'Asia/Shanghai'
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US.UTF-8
ENV LC_ALL en_US.UTF-8
WORKDIR /tmp

RUN DEBIAN_FRONTEND=noninteractive apt-get -yqq install ca-certificates

RUN mkdir -p /var/log/amzimg3

COPY builds/amzimg3 /usr/local/bin/amzimg3
COPY etc/allow_sources.txt /opt/amzimg3/etc/allow_sources.txt

RUN ln -sf /dev/stdout /var/log/amzimg3/stdout.log
RUN ln -sf /dev/stderr /var/log/amzimg3/error.log

VOLUME ["/var/data","/var/log/amzimg3","/opt/amzimg3/etc"]

EXPOSE 8085

ENTRYPOINT ["/usr/local/bin/amzimg3"]
