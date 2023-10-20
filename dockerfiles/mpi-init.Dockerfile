FROM cgr.dev/dominodatalab.com/wolfi-base:1@sha256:0d8fcaa8f8424a0819cf6fb6418bcc82ecfb67ecec6671729f62e725ff27db1d

ARG INSTALL_DIR=/opt/domino/mpi-cluster
ARG INSTALL_BIN=${INSTALL_DIR}/bin

ADD mpi-worker-start.sh ${INSTALL_BIN}/mpi-worker-start.sh

RUN \
    apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache openssh && \
	chmod 755 ${INSTALL_BIN}/mpi-worker-start.sh
