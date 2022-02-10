#!/bin/bash

set -o nounset
set -o errexit

CONFIG_DIR="/opt/domino/etc"

/usr/bin/ssh-keygen -f "$CONFIG_DIR/ssh_host_key" -N '' -t ed25519
chmod 400 "$CONFIG_DIR/ssh_host_key"

/usr/sbin/sshd -f "$CONFIG_DIR/sshd_config" -o "Port $RSYNC_PORT" -De 2>&1 | \
	grep -v 'kex_exchange_identification'