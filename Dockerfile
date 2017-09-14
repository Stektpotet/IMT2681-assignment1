FROM alpine:latest

MAINTAINER Halvor Smedås <stektpotet@gmail.com>

WORKDIR "/opt"

ADD .docker_build/imt2681-assignment1 /opt/bin/imt2681-assignment1
ADD ./templates /opt/templates
ADD ./static /opt/static

CMD ["/opt/bin/imt2681-assignment1"]

