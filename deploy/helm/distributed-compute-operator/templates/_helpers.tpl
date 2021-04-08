{{/*
Return the proper image name
*/}}
{{- define "dco.image" -}}
{{- $imageRoot := .Values.image -}}
{{- $_ := set $imageRoot "tag" (.Values.image.tag | default .Chart.AppVersion) -}}
{{- include "common.images.image" (dict "imageRoot" $imageRoot "global" $) -}}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "dco.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
{{- default (include "common.names.fullname" .) .Values.serviceAccount.name -}}
{{- else -}}
{{- default "default" .Values.serviceAccount.name -}}
{{- end -}}
{{- end -}}

{{/*
Webhook service name
*/}}
{{- define "dco.webhook.service" -}}
{{ include "common.names.fullname" . }}-webhook-server
{{- end -}}

{{/*
Webhook certificate CA name
*/}}
{{- define "dco.webhook.issuer" -}}
{{ include "common.names.fullname" . }}-selfsigned-issuer
{{- end -}}

{{/*
Webhook certificate name
*/}}
{{- define "dco.webhook.certificate" -}}
{{ include "common.names.fullname" . }}-webhook
{{- end -}}

{{/*
Webhook certificate secret name
*/}}
{{- define "dco.webhook.secret" -}}
{{ include "common.names.fullname" . }}-webhook-cert
{{- end -}}

{{/*
Webhook certificate injection annotation
*/}}
{{- define "dco.webhook.annotation" -}}
cert-manager.io/inject-ca-from: {{ .Release.Namespace }}/{{ include "dco.webhook.certificate" . }}
{{- end -}}

{{/*
Returns a name suitable for all manager RBAC objects
*/}}
{{- define "dco.rbac.managerName" -}}
dominodatalab:operator:{{ include "common.names.fullname" . }}:manager
{{- end -}}

{{/*
Returns a name suitable for all hook RBAC objects
*/}}
{{- define "dco.rbac.hookName" -}}
dominodatalab:operator:{{ include "common.names.fullname" . }}:hook
{{- end -}}

{{/*
Returns the appropriate apiVersion for admission registration resources
*/}}
{{- define "dco.capabilities.admissionregistration.apiVersion" -}}
{{- if semverCompare "<1.16-0" (include "common.capabilities.kubeVersion" .) -}}
{{- print "admissionregistration.k8s.io/v1beta1" -}}
{{- else -}}
{{- print "admissionregistration.k8s.io/v1" -}}
{{- end -}}
{{- end -}}
