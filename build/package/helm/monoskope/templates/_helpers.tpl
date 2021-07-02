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
app.kubernetes.io/name: {{ (include "monoskope.name" .) }}
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
NOTE: This utility template is needed until https://git.io/JvuGN is resolved.
Call a template from the context of a subchart.
Usage:
  {{ include "call-nested" (list . "<subchart_name>" "<subchart_template_name>") }}
*/}}
{{- define "call-nested" }}
{{- $dot := index . 0 }}
{{- $subchart := index . 1 | splitList "." }}
{{- $template := index . 2 }}
{{- $values := $dot.Values }}
{{- range $subchart }}
{{- $values = index $values . }}
{{- end }}
{{- include $template (dict "Chart" (dict "Name" (last $subchart)) "Values" $values "Release" $dot.Release "Capabilities" $dot.Capabilities) }}
{{- end }}

{{- define "monoskope.trustAnchorSecretName" -}}
{{- printf "%s-trust-anchor" (include "monoskope.fullname" .) }}
{{- end }}

{{- define "monoskope.tlsSecretName" -}}
{{- printf "%s-tls-cert" (include "monoskope.fullname" .) }}
{{- end }}

{{- define "monoskope.mtlsSecretName" -}}
{{- printf "%s-mtls-cert" (include "monoskope.fullname" .) }}
{{- end }}

{{- define "monoskope.identityCAName" -}}
{{- printf "%s-identity" (include "monoskope.fullname" .) }}
{{- end }}

{{- define "monoskope.domain" -}}
{{- required "a value for .Values.hosting.domain has to be provided" .Values.hosting.domain }}
{{- end }}

{{- define "monoskope.mtlsDomain" -}}
{{- printf "mapi.%s" .Values.hosting.domain }}
{{- end }}

{{- define "monoskope.tlsDomain" -}}
{{- printf "api.%s" .Values.hosting.domain }}
{{- end }}
