#!/bin/bash

set -o nounset
set -o errexit

/usr/bin/rsync \
	--daemon \
	--no-detach \
	--verbose \
       	--config="$DOMINO_ETC/$RSYNC_CONFIG_FILE" \
	--port=$RSYNC_PORT
