FROM ubuntu:14.04
RUN mkdir /usr/mydocker
WORKDIR /usr/mydocker
ADD my-docker /usr/mydocker/
COPY busybox.tar /root/