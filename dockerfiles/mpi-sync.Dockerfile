FROM cgr.dev/dominodatalab.com/chainguard-base@sha256:c14b2aaf63b842a3e65f9af82f1d9dcfa22907c07bbc41f9bdd733a1566dbb36

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
ARG RSYNC_MNT_DIR="/mnt"

ARG ALLENV="\$RSYNC_RUN_DIR,\$DOMINO_ETC,\$RSYNC_CONFIG_FILE"

WORKDIR /root

RUN \
    apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache rsync gettext procps && \
	mkdir -p \
		"$DOMINO_DIR" \
		"$DOMINO_BIN" \
		"$DOMINO_ETC" \
		"$RSYNC_RUN_DIR" \
		"$RSYNC_MNT_DIR"

ADD $RSYNC_START_SCRIPT $RSYNC_CONFIG_FILE ./

RUN \
	addgroup -g $DOMINO_GID -S $DOMINO_GROUP && \
	adduser -u $DOMINO_UID -G $DOMINO_GROUP -D -s /bin/sh $DOMINO_USER && \
	envsubst "$ALLENV" < "$RSYNC_START_SCRIPT" > "$DOMINO_BIN/$RSYNC_START_SCRIPT" && \
	envsubst "$ALLENV" < "$RSYNC_CONFIG_FILE" > "$DOMINO_ETC/$RSYNC_CONFIG_FILE" && \
	chown -R $DOMINO_USER:$DOMINO_GROUP "$RSYNC_RUN_DIR" && \
	chown -R $DOMINO_USER:$DOMINO_GROUP "$DOMINO_DIR" && \
	chmod 755 "$DOMINO_BIN/$RSYNC_START_SCRIPT" && \
	chmod 644 "$DOMINO_ETC/$RSYNC_CONFIG_FILE"

# For testing -- to be removed
RUN \
	chown -R $DOMINO_USER:$DOMINO_GROUP /mnt
