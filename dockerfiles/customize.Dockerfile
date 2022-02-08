# A specific version of the Linux OS here is very important, because 
# it defines versions of core libraries (libc etc) the compiled binaries
# will be linked against.
FROM ubuntu:18.04

ARG RSYNC_VERSION=3.2.3
ARG RSYNC_URL=https://download.samba.org/pub/rsync/src/rsync-${RSYNC_VERSION}.tar.gz
ARG RSYNC_SIG_URL=https://download.samba.org/pub/rsync/src/rsync-${RSYNC_VERSION}.tar.gz.asc

ARG OPENSSH_VERSION=8.8p1
ARG OPENSSH_URL=https://mirrors.mit.edu/pub/OpenBSD/OpenSSH/portable/openssh-${OPENSSH_VERSION}.tar.gz
ARG OPENSSH_SIG_URL=https://mirrors.mit.edu/pub/OpenBSD/OpenSSH/portable/openssh-${OPENSSH_VERSION}.tar.gz.asc

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
	gpg --import -q openssh.gpgkey > /dev/null && \
	rm -f *.gpgkey

# Download and compile rsync 
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

# Download an compile openssh
RUN \
	# Newer versions of openssh include a mandatory privilege separation mechanism
	# that requires a special user to be available in the system. Although this
	# image does not execute sshd, such a user must exist for proper deployment.
	useradd -g 65534 -d /var/empty -s /bin/false sshd && \
	curl -o openssh-src.tgz -LSsf ${OPENSSH_URL} && \
	curl -o openssh-src.sig -LSsf ${OPENSSH_SIG_URL} && \
	gpg --trust-model always -q --verify openssh-src.sig openssh-src.tgz && \
	tar -xf openssh-src.tgz --no-same-permissions && \
	cd openssh-${OPENSSH_VERSION} && \
	./configure \
		--prefix=${INSTALL_DIR} \
		--without-zlib \
		--without-openssl && \
	make && \
	make install && \
	cd -

ADD mpi-worker-start.sh ${INSTALL_BIN}

# Create a tarball containing all the necessary stuff
RUN \
	chmod 755 ${INSTALL_BIN}/mpi-worker-start.sh && \
	tar -czf worker-utils.tgz \
		${INSTALL_DIR}/bin \
		${INSTALL_DIR}/etc \
		${INSTALL_DIR}/libexec \
		${INSTALL_DIR}/sbin

# The final image only contains built artifacts.
FROM debian:11-slim
WORKDIR /root
COPY --from=0 /root/worker-utils.tgz ./
CMD tar -C / -xf /root/worker-utils.tgz
