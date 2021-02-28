{{/*
Expand the name of the chart.
*/}}
{{- define "distributed-compute-operator.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "distributed-compute-operator.fullname" -}}
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
{{- define "distributed-compute-operator.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "distributed-compute-operator.labels" -}}
helm.sh/chart: {{ include "distributed-compute-operator.chart" . }}
{{ include "distributed-compute-operator.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "distributed-compute-operator.selectorLabels" -}}
app.kubernetes.io/name: {{ include "distributed-compute-operator.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "distributed-compute-operator.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "distributed-compute-operator.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Webhook service name
*/}}
{{- define "distributed-compute-operator.webhook.service" -}}
{{ include "distributed-compute-operator.fullname" . }}-webhook-server
{{- end }}

{{/*
Webhook certificate CA name
*/}}
{{- define "distributed-compute-operator.webhook.issuer" -}}
{{ include "distributed-compute-operator.fullname" . }}-selfsigned-issuer
{{- end }}

{{/*
Webhook certificate name
*/}}
{{- define "distributed-compute-operator.webhook.certificate" -}}
{{ include "distributed-compute-operator.fullname" . }}-webhook
{{- end }}

{{/*
Webhook certificate secret name
*/}}
{{- define "distributed-compute-operator.webhook.secret" -}}
{{ include "distributed-compute-operator.fullname" . }}-webhook-cert
{{- end }}
