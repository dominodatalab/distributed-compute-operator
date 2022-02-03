# A specific version of the Linux OS here is very important, because 
# it defines versions of core libraries (libc etc) the compiled binaries
# will be linked against.
#FROM debian:9.13
FROM ubuntu:18.04

# Locations of the source code for the utilities 
ARG RSYNC_VERSION=3.2.3
ARG RSYNC_URL=https://download.samba.org/pub/rsync/src/rsync-${RSYNC_VERSION}.tar.gz
ARG RSYNC_SIG_URL=https://download.samba.org/pub/rsync/src/rsync-${RSYNC_VERSION}.tar.gz.asc

ARG INSTALL_DIR=/opt/domino
ARG INSTALL_BIN=${INSTALL_DIR}/bin

WORKDIR /root

ADD *.gpgkey ./

# Install common dependencies for the compiler and setting things up
RUN \
	apt-get update && \
	apt-get -y install \
		build-essential \
		curl \
		gnupg && \
	 mkdir -p \
		${INSTALL_DIR} \
		${INSTALL_BIN} && \
	gpg --import -q rsync.gpgkey > /dev/null && \
#	gpg --import -q curl.gpgkey > /dev/null && \
	rm -f *.gpgkey

# Download and compile sudo
RUN \
	curl -o rsync-src.tgz -LSsf ${RSYNC_URL} && \
	curl -o rsync-src.sig -LSsf ${RSYNC_SIG_URL} && \
	gpg --trust-model always -q --verify rsync-src.sig rsync-src.tgz && \
	tar -xf rsync-src.tgz --no-same-permissions && \
	cd rsync-${RSYNC_VERSION} && \
	./configure \
		--prefix=${INSTALL_DIR} \
		--disable-openssl \
		--disable-xxhash \
		--disable-zstd \
		--disable-lz4 && \
	make && \
	make install && \
	cd - 

CMD /bin/bash
