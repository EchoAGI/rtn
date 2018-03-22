# Channel server Docker builder
#
# This Dockerfile creates a container which builds Channel Server as found in the
# current folder, and creates a tarball which can be piped into another Docker
# container for creating minimal sized Docker containers.
#
# First create the builder image:
#
#   ```
#   docker build -t channel-server-builder -f Dockerfile.build .
#   ```
# Next run the builder container, piping its output into the creation of the
# runner container. This creates a minimal size Docker image which can be used
# to run Channel Server in production.
#
#   ```
#   docker run --rm channel-server-builder | docker build -t channel-server -f Dockerfile.run -
#   ```

#FROM ubuntu:xenial
FROM golang:latest
MAINTAINER edison <52388483@qq.com>

# Set locale.
#RUN locale-gen --no-purge en_US.UTF-8
#ENV LC_ALL en_US.UTF-8

ENV DEBIAN_FRONTEND noninteractive

# Base build dependencies.
RUN apt-get update && apt-get install -qy \
	nodejs \
	build-essential \
	git \
	automake \
	autoconf

# Add and build Channel server.
ADD . /srv/channel-server
WORKDIR /srv/channel-server

RUN git clone http://github.com/golang/net.git /go/src/golang.org/x/net \
    && git clone http://github.com/golang/sys.git /go/src/golang.org/x/sys \
	&& git clone http://github.com/golang/crypto.git /go/src/golang.org/x/crypto

RUN mkdir -p /usr/share/gocode/src
RUN ./autogen.sh && \
	./configure && \
	make pristine && \
	make get && \
	make tarball
RUN rm /srv/channel-server/dist_*/*.tar.gz
RUN mv /srv/channel-server/dist_*/channel-server-* /srv/channel-server/dist

# Add gear required by Dockerfile.run.
COPY Dockerfile.run /
COPY scripts/docker_entrypoint.sh /

# Running this image produces a tarball suitable to be piped into another Docker build command.
CMD tar -cf - -C / Dockerfile.run docker_entrypoint.sh /srv/channel-server/dist
