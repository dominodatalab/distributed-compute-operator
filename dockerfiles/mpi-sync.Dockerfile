FROM quay.io/domino/debian:10.11-20220621-1030

ARG DOMINO_UID=12574
ARG DOMINO_USER=domino
ARG DOMINO_GID=12574
ARG DOMINO_GROUP=domino

ARG DOMINO_DIR=/opt/domino/rsync
ARG DOMINO_BIN=$DOMINO_DIR/bin
ARG DOMINO_ETC=$DOMINO_DIR/etc

ARG RSYNC_RUN_DIR=/run/rsyncd-${DOMINO_USER}
ARG RSYNC_CONFIG_FILE=rsyncd.conf
ARG RSYNC_START_SCRIPT=rsync-start.sh

ARG ALLENV="\$RSYNC_RUN_DIR,\$DOMINO_ETC,\$RSYNC_CONFIG_FILE"

WORKDIR /root

RUN \
	apt-get update && \
	apt-get -y install \
		rsync \
		gettext-base \
		procps && \
	rm -rf /var/lib/apt/lists/* && \
	mkdir -p \
		"$DOMINO_DIR" \
		"$DOMINO_BIN" \
		"$DOMINO_ETC" \
		"$RSYNC_RUN_DIR"

ADD $RSYNC_START_SCRIPT $RSYNC_CONFIG_FILE ./

RUN \
	groupadd -g $DOMINO_GID $DOMINO_GROUP && \
	useradd -u $DOMINO_UID -g $DOMINO_GID -mN -s /bin/bash $DOMINO_USER && \
	envsubst "$ALLENV" < "$RSYNC_START_SCRIPT" > "$DOMINO_BIN/$RSYNC_START_SCRIPT" && \
	envsubst "$ALLENV" < "$RSYNC_CONFIG_FILE" > "$DOMINO_ETC/$RSYNC_CONFIG_FILE" && \
	chown -R $DOMINO_USER:$DOMINO_GROUP "$RSYNC_RUN_DIR" && \
	chown -R $DOMINO_USER:$DOMINO_GROUP "$DOMINO_DIR" && \
	chmod 755 "$DOMINO_BIN/$RSYNC_START_SCRIPT" && \
	chmod 644 "$DOMINO_ETC/$RSYNC_CONFIG_FILE"

# For testing -- to be removed
RUN \
	chown -R $DOMINO_USER:$DOMINO_GROUP /mnt
