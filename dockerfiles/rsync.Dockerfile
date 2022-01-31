FROM debian:11-slim

ARG DOMINO_UID=12574
ARG DOMINO_USER=domino
ARG DOMINO_GID=12574
ARG DOMINO_GROUP=domino

ARG RSYNC_SSH_PORT=2223

ARG DOMINO_DIR=/opt/domino
ARG DOMINO_SSH_DIR=${DOMINO_DIR}/etc/ssh

ARG SSHD_CONFIG=${DOMINO_SSH_DIR}/sshd_config
ARG AUTHORIZED_KEYS_PATH=/etc/mpi/authorized_keys

RUN \
	apt-get update && \
	apt-get -y install \
		openssh-server \
		rsync && \
    rm -rf /var/lib/apt/lists/* && \
    rm -rf /etc/ssh/ssh_host*

RUN \
    groupadd -g ${DOMINO_GID} ${DOMINO_GROUP} && \
	useradd -u ${DOMINO_UID} -g ${DOMINO_GID} -mN -s /bin/bash ${DOMINO_USER}

RUN \
    mkdir -p ${DOMINO_DIR} ${DOMINO_SSH_DIR} && \
    rm -f ${SSHD_CONFIG} && \
    echo "HostKey \"${DOMINO_SSH_DIR}/ssh_host_key\"" >> ${SSHD_CONFIG} && \
    echo "AuthorizedKeysFile \"${AUTHORIZED_KEYS_PATH}\"" >> ${SSHD_CONFIG} && \
    echo "PidFile \"/tmp/domino-sshd.pid\"" >> ${SSHD_CONFIG} && \
    echo "AllowUsers ${DOMINO_USER}" >> ${SSHD_CONFIG} && \
    chmod 444 ${SSHD_CONFIG} && \
    chown -R ${DOMINO_USER}:${DOMINO_GROUP} ${DOMINO_DIR}
