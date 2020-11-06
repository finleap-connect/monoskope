{{/*
Expand the name of the chart.
*/}}
{{- define "monoskope.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "monoskope.fullname" -}}
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
{{- define "monoskope.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "monoskope.labels" -}}
helm.sh/chart: {{ include "monoskope.chart" . }}
{{ include "monoskope.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "monoskope.selectorLabels" -}}
app.kubernetes.io/name: {{ include "monoskope.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "monoskope.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "monoskope.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Expand the name of the chart.
*/}}
{{- define "monoskope.name.gateway" -}}
{{- default (printf "%s-%s" (include "monoskope.name" .) "gateway") .Values.gateway.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "monoskope.fullname.gateway" -}}
{{- default (printf "%s-%s" (include "monoskope.fullname" .) "gateway")  .Values.gateway.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "monoskope.labels.gateway" -}}
helm.sh/chart: {{ include "monoskope.chart" . }}
{{ include "monoskope.selectorLabels.gateway" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "monoskope.selectorLabels.gateway" -}}
app.kubernetes.io/name: {{ include "monoskope.name.gateway" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Expand the name of the chart.
*/}}
{{- define "monoskope.name.ingress" -}}
{{- (printf "%s-%s" (include "monoskope.name" .) "ingress") | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "monoskope.fullname.ingress" -}}
{{- (printf "%s-%s" (include "monoskope.fullname" .) "ingress")  | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "monoskope.labels.ingress" -}}
helm.sh/chart: {{ include "monoskope.chart" . }}
{{ include "monoskope.selectorLabels.ingress" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "monoskope.selectorLabels.ingress" -}}
app.kubernetes.io/name: {{ include "monoskope.name.ingress" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}