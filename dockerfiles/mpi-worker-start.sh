#!/bin/bash

set -o nounset
set -o errexit

INSTALL_DIR="/opt/domino/mpi-cluster"
SSH_USER="sshd"
SSH_RUN_DIR="/run/sshd-${DOMINO_USER}"

mkdir -p "$SSH_RUN_DIR"
chmod 777 "$SSH_RUN_DIR"

if ! id $SSH_USER >/dev/null 2>&1; then
	useradd -g 65534 -mN -s "/usr/sbin/nologin" -d "$SSH_RUN_DIR" $SSH_USER
fi

if ! cut -d: -f3 < /etc/group | grep "^${DOMINO_GID}$" >/dev/null 2>&1; then
	groupadd -g $DOMINO_GID $DOMINO_GROUP
fi
if ! id $DOMINO_UID >/dev/null 2>&1; then
	useradd -u $DOMINO_UID -g $DOMINO_GID -mN -s /bin/bash -d "$DOMINO_HOME_DIR" $DOMINO_USER
else
    # Change username of user with matching userid if needed
	EXISTING_USER=$(id -nu $DOMINO_UID)
	if [ "$EXISTING_USER" != "$DOMINO_USER" ]; then
		usermod -l $DOMINO_USER $EXISTING_USER
	fi

	# Change home directory (idempotent)
	usermod -d "$DOMINO_HOME_DIR" $DOMINO_USER

    # Add to domino group (idempotent)
    usermod -a -G $DOMINO_GROUP $DOMINO_USER
fi

# Add the new domino user to the non-root groups of the current container user
for gid in `id -G`; do
  if [ $gid != 0 ]; then
    # Add user to a new/existing group with desired id.
    group_name=$(cut -d: -f1 < <(getent group $gid))
    if [ -z $group_name ]; then
        group_name="group-$gid"
	    groupadd -g $gid $group_name
    fi
    usermod -a -G $group_name $DOMINO_USER
  fi
done

CONFIG_DIR="$INSTALL_DIR/etc"
mkdir -p "$CONFIG_DIR"

rm -f "$CONFIG_DIR/ssh_host_*"
"$INSTALL_DIR/bin/ssh-keygen" -f "$CONFIG_DIR/ssh_host_key" -N '' -t ed25519
chmod 400 "$CONFIG_DIR/ssh_host_key"
chown $DOMINO_UID:$DOMINO_GID "$CONFIG_DIR/ssh_host_key"

cat << EOF > "$CONFIG_DIR/sshd_config"
Port $DOMINO_SSH_PORT
HostKey "$CONFIG_DIR/ssh_host_key"
AuthorizedKeysFile "$DOMINO_KEYS_PATH"
PidFile "$SSH_RUN_DIR/sshd.pid"
AllowUsers $DOMINO_USER
EOF
chmod 444 "$CONFIG_DIR/sshd_config"

su -c "$INSTALL_DIR/sbin/sshd -f \"$CONFIG_DIR/sshd_config\" -De" - $DOMINO_USER 
