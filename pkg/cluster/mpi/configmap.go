package mpi

import (
	"fmt"
	"path/filepath"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	dcv1alpha1 "github.com/dominodatalab/distributed-compute-operator/api/v1alpha1"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/actions"
	"github.com/dominodatalab/distributed-compute-operator/pkg/controller/core"
)

const launchScriptTemplate = `#!/bin/bash

set -o nounset 
set -o errexit

SSH_PORT=%[1]d
DOMINO_UID=%[2]d
DOMINO_GID=%[3]d
DOMINO_USER=%[4]s
DOMINO_GROUP=%[5]s
AUTHORIZED_KEYS_PATH=%[6]s

if ! id $DOMINO_UID >/dev/null 2>&1; then
    groupadd -g $DOMINO_GID $DOMINO_GROUP
    useradd -u $DOMINO_UID -g $DOMINO_GID -mN -s /bin/bash $DOMINO_USER
fi
if [ "$(id -nu $DOMINO_UID)" != "$DOMINO_USER" ]; then
    echo >&2 "User name mismatch"
    exit 1
fi

DOMINO_HOME=$(eval echo "~$DOMINO_USER")
SSH_DIR="$DOMINO_HOME/.ssh"
mkdir -p "$SSH_DIR"

ssh-keygen -f "$SSH_DIR/ssh_host_key" -N '' -t ecdsa
chmod 400 "$SSH_DIR/ssh_host_key"

cat << EOF > "$SSH_DIR/sshd_config"
Port $SSH_PORT
HostKey "$SSH_DIR/ssh_host_key"
AuthorizedKeysFile "$AUTHORIZED_KEYS_PATH"
PidFile "$SSH_DIR/sshd.pid"
AllowUsers $DOMINO_USER
EOF
chmod 444 "$SSH_DIR/sshd_config"

chown -R $DOMINO_UID:$DOMINO_GID "$SSH_DIR"
chmod 755 "$SSH_DIR"

su -c "/usr/sbin/sshd -f \"$SSH_DIR/sshd_config\" -De" - $DOMINO_USER`

func ConfigMap() core.OwnedComponent {
	return &configMapComponent{}
}

type configMapComponent struct{}

func (c configMapComponent) Reconcile(ctx *core.Context) (ctrl.Result, error) {
	cr := objToMPICluster(ctx.Object)

	launchScriptName := filepath.Base(launchScriptPath)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName(cr),
			Namespace: cr.Namespace,
			Labels:    meta.StandardLabels(cr),
		},
		Data: map[string]string{
			hostFileName:     buildHostFile(cr),
			launchScriptName: buildLaunchScript(cr),
		},
	}

	err := actions.CreateOrUpdateOwnedResource(ctx, cr, cm)
	if err != nil {
		err = fmt.Errorf("cannot reconcile configmap: %w", err)
	}

	return ctrl.Result{}, err
}

func (c configMapComponent) Kind() client.Object {
	return &corev1.ConfigMap{}
}

func buildHostFile(cr *dcv1alpha1.MPICluster) string {
	svcName := serviceName(cr)
	workerName := workerStatefulSetName(cr)
	workerReplicas := *cr.Spec.Worker.Replicas

	var builder strings.Builder
	for idx := 0; idx < int(workerReplicas); idx++ {
		entry := fmt.Sprintf("%s-%d.%s\n", workerName, idx, svcName)
		builder.WriteString(entry)
	}

	return builder.String()
}

func buildLaunchScript(cr *dcv1alpha1.MPICluster) string {
	userId := int64(defaultUserID)
	if cr.Spec.Worker.UserId != nil {
		userId = *cr.Spec.Worker.UserId
	}
	userName := defaultUserName
	if cr.Spec.Worker.UserName != "" {
		userName = cr.Spec.Worker.UserName
	}
	groupId := int64(defaultGroupID)
	if cr.Spec.Worker.GroupId != nil {
		groupId = *cr.Spec.Worker.GroupId
	}
	groupName := defaultGroupName
	if cr.Spec.Worker.GroupName != "" {
		groupName = cr.Spec.Worker.GroupName
	}

	return fmt.Sprintf(launchScriptTemplate,
		sshdPort,           // 1 int
		userId,             // 2 int
		groupId,            // 3 int
		userName,           // 4 string
		groupName,          // 5 string
		authorizedKeysPath, // 6 string
	)
}
