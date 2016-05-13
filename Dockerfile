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

COPY builds/amzimg3 /usr/local/bin/amzimg3
COPY etc/allow_sources.txt /opt/amzimg3/etc/allow_sources.txt

VOLUME ["/var/data","/opt/amzimg3/etc"]

EXPOSE 8085

ENTRYPOINT ["/usr/local/bin/amzimg3"]