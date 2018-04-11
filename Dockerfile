FROM alpine:latest

MAINTAINER weidewang dewang.wei@gmail.com

RUN apk add --update --no-cache ca-certificates tzdata

ENV TZ 'Asia/Shanghai'
ENV LANG en_US.UTF-8
ENV LANGUAGE en_US.UTF-8
ENV LC_ALL en_US.UTF-8
WORKDIR /tmp


RUN mkdir -p /var/log/amzimg3

COPY builds/amzimg3 /usr/local/bin/amzimg3
COPY etc/allow_sources.txt /opt/amzimg3/etc/allow_sources.txt

VOLUME ["/var/data","/var/log/amzimg3","/opt/amzimg3/etc"]

EXPOSE 8085

ENTRYPOINT ["/usr/local/bin/amzimg3"]
