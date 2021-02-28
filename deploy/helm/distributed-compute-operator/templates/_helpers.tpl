{{/*
Expand the name of the chart.
*/}}
{{- define "dco.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "dco.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "dco.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "dco.labels" -}}
helm.sh/chart: {{ include "dco.chart" . }}
{{ include "dco.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "dco.selectorLabels" -}}
app.kubernetes.io/name: {{ include "dco.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "dco.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "dco.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Webhook service name
*/}}
{{- define "dco.webhook.service" -}}
{{ include "dco.fullname" . }}-webhook-server
{{- end }}

{{/*
Webhook certificate CA name
*/}}
{{- define "dco.webhook.issuer" -}}
{{ include "dco.fullname" . }}-selfsigned-issuer
{{- end }}

{{/*
Webhook certificate name
*/}}
{{- define "dco.webhook.certificate" -}}
{{ include "dco.fullname" . }}-webhook
{{- end }}

{{/*
Webhook certificate secret name
*/}}
{{- define "dco.webhook.secret" -}}
{{ include "dco.fullname" . }}-webhook-cert
{{- end }}
