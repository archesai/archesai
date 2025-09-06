{{/*
Expand the name of the chart.
*/}}
{{- define "arches.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "arches.fullname" -}}
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
{{- define "arches.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "arches.labels" -}}
helm.sh/chart: {{ include "arches.chart" . }}
{{ include "arches.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "arches.selectorLabels" -}}
app.kubernetes.io/name: {{ include "arches.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "arches.serviceAccountName" -}}
{{- if .Values.infrastructure.serviceAccount.create }}
{{- default (include "arches.fullname" .) .Values.infrastructure.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.infrastructure.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Generate the config YAML for the application
This takes the config section from values and converts it to YAML
*/}}
{{- define "arches.config" -}}
{{- $config := deepCopy .Values }}
{{- /* Update service references to match Kubernetes service names */ -}}
{{- if eq $config.redis.mode "managed" }}
{{- $_ := set $config.redis "host" (printf "%s-redis" (include "arches.fullname" .)) }}
{{- end }}
{{- if eq $config.storage.mode "managed" }}
{{- $_ := set $config.storage "endpoint" (printf "http://%s-minio:9000" (include "arches.fullname" .)) }}
{{- end }}
{{- if eq $config.intelligence.scraper.mode "managed" }}
{{- $_ := set $config.intelligence.scraper "endpoint" (printf "http://%s-scraper:8080" (include "arches.fullname" .)) }}
{{- end }}
{{- if eq $config.intelligence.unstructured.mode "managed" }}
{{- $_ := set $config.intelligence.unstructured "endpoint" (printf "http://%s-unstructured:8000" (include "arches.fullname" .)) }}
{{- end }}
{{- if eq $config.monitoring.loki.mode "managed" }}
{{- $_ := set $config.monitoring.loki "host" (printf "http://%s-loki:3100" (include "arches.fullname" .)) }}
{{- end }}
{{- $config | toYaml }}
{{- end }}