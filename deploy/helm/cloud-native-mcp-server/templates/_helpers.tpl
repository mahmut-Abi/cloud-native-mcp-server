{{/*
Expand the name of the chart.
*/}}
{{- define "cloud-native-mcp-server.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "cloud-native-mcp-server.fullname" -}}
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
{{- define "cloud-native-mcp-server.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cloud-native-mcp-server.labels" -}}
helm.sh/chart: {{ include "cloud-native-mcp-server.chart" . }}
{{ include "cloud-native-mcp-server.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cloud-native-mcp-server.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cloud-native-mcp-server.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "cloud-native-mcp-server.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "cloud-native-mcp-server.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Create the image name
*/}}
{{- define "cloud-native-mcp-server.image" -}}
{{- $registry := .Values.global.imageRegistry | default .Values.image.registry -}}
{{- $repository := .Values.image.repository -}}
{{- $tag := .Values.image.tag | default .Chart.AppVersion -}}
{{- if .Values.global.imageRegistry -}}
{{- printf "%s/%s:%s" $registry $repository $tag -}}
{{- else -}}
{{- printf "%s/%s:%s" $registry $repository $tag -}}
{{- end -}}
{{- end }}

{{/*
Return the proper Docker Image Registry Secret Names
*/}}
{{- define "cloud-native-mcp-server.imagePullSecrets" -}}
{{- include "common.images.pullSecrets" (dict "images" (list .Values.image) "global" .Values.global) -}}
{{- end }}

{{/*
Create the configmap checksum
*/}}
{{- define "cloud-native-mcp-server.configChecksum" -}}
{{- $config := include "cloud-native-mcp-server.config" . | fromYaml }}
{{- $config | toYaml | sha256sum }}
{{- end }}

{{/*
Create the configuration
*/}}
{{- define "cloud-native-mcp-server.config" -}}
kubernetes:
  kubeconfig: {{ .Values.config.kubernetes.kubeconfig | quote }}
  timeoutSec: {{ .Values.config.kubernetes.timeoutSec }}
  qps: {{ .Values.config.kubernetes.qps }}
  burst: {{ .Values.config.kubernetes.burst }}

server:
  mode: {{ .Values.config.server.mode | quote }}
  addr: {{ .Values.config.server.addr | quote }}
  readTimeoutSec: {{ .Values.config.server.readTimeoutSec }}
  writeTimeoutSec: {{ .Values.config.server.writeTimeoutSec }}
  idleTimeoutSec: {{ .Values.config.server.idleTimeoutSec }}
  ssePaths:
    kubernetes: {{ .Values.config.server.ssePaths.kubernetes | quote }}
    grafana: {{ .Values.config.server.ssePaths.grafana | quote }}
    prometheus: {{ .Values.config.server.ssePaths.prometheus | quote }}
    loki: {{ .Values.config.server.ssePaths.loki | quote }}
    kibana: {{ .Values.config.server.ssePaths.kibana | quote }}
    helm: {{ .Values.config.server.ssePaths.helm | quote }}
    elasticsearch: {{ .Values.config.server.ssePaths.elasticsearch | quote }}
    alertmanager: {{ .Values.config.server.ssePaths.alertmanager | quote }}
    jaeger: {{ .Values.config.server.ssePaths.jaeger | quote }}
    opentelemetry: {{ .Values.config.server.ssePaths.opentelemetry | quote }}
    utilities: {{ .Values.config.server.ssePaths.utilities | quote }}
    aggregate: {{ .Values.config.server.ssePaths.aggregate | quote }}
  streamableHttpPaths:
    kubernetes: {{ .Values.config.server.streamableHttpPaths.kubernetes | quote }}
    grafana: {{ .Values.config.server.streamableHttpPaths.grafana | quote }}
    prometheus: {{ .Values.config.server.streamableHttpPaths.prometheus | quote }}
    loki: {{ .Values.config.server.streamableHttpPaths.loki | quote }}
    kibana: {{ .Values.config.server.streamableHttpPaths.kibana | quote }}
    helm: {{ .Values.config.server.streamableHttpPaths.helm | quote }}
    elasticsearch: {{ .Values.config.server.streamableHttpPaths.elasticsearch | quote }}
    alertmanager: {{ .Values.config.server.streamableHttpPaths.alertmanager | quote }}
    jaeger: {{ .Values.config.server.streamableHttpPaths.jaeger | quote }}
    opentelemetry: {{ .Values.config.server.streamableHttpPaths.opentelemetry | quote }}
    utilities: {{ .Values.config.server.streamableHttpPaths.utilities | quote }}
    aggregate: {{ .Values.config.server.streamableHttpPaths.aggregate | quote }}
  cors:
    allowedOrigins:
      {{- range .Values.config.server.cors.allowedOrigins }}
      - {{ . | quote }}
      {{- end }}
    allowedMethods:
      {{- range .Values.config.server.cors.allowedMethods }}
      - {{ . | quote }}
      {{- end }}
    allowedHeaders:
      {{- range .Values.config.server.cors.allowedHeaders }}
      - {{ . | quote }}
      {{- end }}
    maxAge: {{ .Values.config.server.cors.maxAge }}

logging:
  level: {{ .Values.config.logging.level | quote }}
  json: {{ .Values.config.logging.json }}

prometheus:
  enabled: {{ .Values.config.prometheus.enabled }}
  address: {{ .Values.config.prometheus.address | quote }}
  bearerToken: {{ .Values.config.prometheus.bearerToken | quote }}
  tlsSkipVerify: {{ .Values.config.prometheus.tlsSkipVerify }}
  timeoutSec: {{ .Values.config.prometheus.timeoutSec }}
  username: {{ .Values.config.prometheus.username | quote }}
  password: {{ .Values.config.prometheus.password | quote }}
  tlsCertFile: {{ .Values.config.prometheus.tlsCertFile | quote }}
  tlsKeyFile: {{ .Values.config.prometheus.tlsKeyFile | quote }}
  tlsCAFile: {{ .Values.config.prometheus.tlsCAFile | quote }}

loki:
  enabled: {{ .Values.config.loki.enabled }}
  address: {{ .Values.config.loki.address | quote }}
  bearerToken: {{ .Values.config.loki.bearerToken | quote }}
  tlsSkipVerify: {{ .Values.config.loki.tlsSkipVerify }}
  timeoutSec: {{ .Values.config.loki.timeoutSec }}
  username: {{ .Values.config.loki.username | quote }}
  password: {{ .Values.config.loki.password | quote }}
  tlsCertFile: {{ .Values.config.loki.tlsCertFile | quote }}
  tlsKeyFile: {{ .Values.config.loki.tlsKeyFile | quote }}
  tlsCAFile: {{ .Values.config.loki.tlsCAFile | quote }}

grafana:
  enabled: {{ .Values.config.grafana.enabled }}
  url: {{ .Values.config.grafana.url | quote }}
  apiKey: {{ .Values.config.grafana.apiKey | quote }}
  username: {{ .Values.config.grafana.username | quote }}
  password: {{ .Values.config.grafana.password | quote }}
  timeoutSec: {{ .Values.config.grafana.timeoutSec }}

kibana:
  enabled: {{ .Values.config.kibana.enabled }}
  url: {{ .Values.config.kibana.url | quote }}
  apiKey: {{ .Values.config.kibana.apiKey | quote }}
  username: {{ .Values.config.kibana.username | quote }}
  password: {{ .Values.config.kibana.password | quote }}
  skipVerify: {{ .Values.config.kibana.skipVerify }}
  space: {{ .Values.config.kibana.space | quote }}
  timeoutSec: {{ .Values.config.kibana.timeoutSec }}

helm:
  enabled: {{ .Values.config.helm.enabled }}
  kubeconfigPath: {{ .Values.config.helm.kubeconfigPath | quote }}
  namespace: {{ .Values.config.helm.namespace | quote }}
  debug: {{ .Values.config.helm.debug }}
  timeoutSec: {{ .Values.config.helm.timeoutSec }}
  maxRetries: {{ .Values.config.helm.maxRetries }}
  httpProxy: {{ .Values.config.helm.httpProxy | quote }}

elasticsearch:
  enabled: {{ .Values.config.elasticsearch.enabled }}
  addresses:
    {{- range .Values.config.elasticsearch.addresses }}
    - {{ . | quote }}
    {{- end }}
  address: {{ .Values.config.elasticsearch.address | quote }}
  username: {{ .Values.config.elasticsearch.username | quote }}
  password: {{ .Values.config.elasticsearch.password | quote }}
  bearerToken: {{ .Values.config.elasticsearch.bearerToken | quote }}
  apiKey: {{ .Values.config.elasticsearch.apiKey | quote }}
  timeoutSec: {{ .Values.config.elasticsearch.timeoutSec }}
  tlsSkipVerify: {{ .Values.config.elasticsearch.tlsSkipVerify }}
  tlsCertFile: {{ .Values.config.elasticsearch.tlsCertFile | quote }}
  tlsKeyFile: {{ .Values.config.elasticsearch.tlsKeyFile | quote }}
  tlsCAFile: {{ .Values.config.elasticsearch.tlsCAFile | quote }}

alertmanager:
  enabled: {{ .Values.config.alertmanager.enabled }}
  address: {{ .Values.config.alertmanager.address | quote }}
  bearerToken: {{ .Values.config.alertmanager.bearerToken | quote }}
  tlsSkipVerify: {{ .Values.config.alertmanager.tlsSkipVerify }}
  timeoutSec: {{ .Values.config.alertmanager.timeoutSec }}
  username: {{ .Values.config.alertmanager.username | quote }}
  password: {{ .Values.config.alertmanager.password | quote }}
  tlsCertFile: {{ .Values.config.alertmanager.tlsCertFile | quote }}
  tlsKeyFile: {{ .Values.config.alertmanager.tlsKeyFile | quote }}
  tlsCAFile: {{ .Values.config.alertmanager.tlsCAFile | quote }}

jaeger:
  enabled: {{ .Values.config.jaeger.enabled }}
  address: {{ .Values.config.jaeger.address | quote }}
  timeoutSec: {{ .Values.config.jaeger.timeoutSec }}

opentelemetry:
  enabled: {{ .Values.config.opentelemetry.enabled }}
  address: {{ .Values.config.opentelemetry.address | quote }}
  bearerToken: {{ .Values.config.opentelemetry.bearerToken | quote }}
  tlsSkipVerify: {{ .Values.config.opentelemetry.tlsSkipVerify }}
  timeoutSec: {{ .Values.config.opentelemetry.timeoutSec }}
  username: {{ .Values.config.opentelemetry.username | quote }}
  password: {{ .Values.config.opentelemetry.password | quote }}
  tlsCertFile: {{ .Values.config.opentelemetry.tlsCertFile | quote }}
  tlsKeyFile: {{ .Values.config.opentelemetry.tlsKeyFile | quote }}
  tlsCAFile: {{ .Values.config.opentelemetry.tlsCAFile | quote }}

auth:
  enabled: {{ .Values.config.auth.enabled }}
  mode: {{ .Values.config.auth.mode | quote }}
  apiKey: {{ .Values.config.auth.apiKey | quote }}
  bearerToken: {{ .Values.config.auth.bearerToken | quote }}
  username: {{ .Values.config.auth.username | quote }}
  password: {{ .Values.config.auth.password | quote }}
  jwtSecret: {{ .Values.config.auth.jwtSecret | quote }}
  jwtAlgorithm: {{ .Values.config.auth.jwtAlgorithm | quote }}

audit:
  enabled: {{ .Values.config.audit.enabled }}
  level: {{ .Values.config.audit.level | quote }}
  storage: {{ .Values.config.audit.storage | quote }}
  format: {{ .Values.config.audit.format | quote }}
  file:
    path: {{ .Values.config.audit.file.path | quote }}
    maxSizeMB: {{ .Values.config.audit.file.maxSizeMB }}
    maxBackups: {{ .Values.config.audit.file.maxBackups }}
    maxAgeDays: {{ .Values.config.audit.file.maxAgeDays }}
    compress: {{ .Values.config.audit.file.compress }}
  database:
    type: {{ .Values.config.audit.database.type | quote }}
    sqlitePath: {{ .Values.config.audit.database.sqlitePath | quote }}
    connectionString: {{ .Values.config.audit.database.connectionString | quote }}
    tableName: {{ .Values.config.audit.database.tableName | quote }}
    poolSize: {{ .Values.config.audit.database.poolSize }}
  fields:
    timestamp: {{ .Values.config.audit.fields.timestamp }}
    clientIP: {{ .Values.config.audit.fields.clientIP }}
    user: {{ .Values.config.audit.fields.user }}
    toolName: {{ .Values.config.audit.fields.toolName }}
    serviceName: {{ .Values.config.audit.fields.serviceName }}
    arguments: {{ .Values.config.audit.fields.arguments }}
    result: {{ .Values.config.audit.fields.result }}
    duration: {{ .Values.config.audit.fields.duration }}
    status: {{ .Values.config.audit.fields.status }}
    error: {{ .Values.config.audit.fields.error }}
  masking:
    enabled: {{ .Values.config.audit.masking.enabled }}
    fields: {{ .Values.config.audit.masking.fields | quote }}
    maskValue: {{ .Values.config.audit.masking.maskValue | quote }}
  sampling:
    enabled: {{ .Values.config.audit.sampling.enabled }}
    rate: {{ .Values.config.audit.sampling.rate }}
  query:
    enabled: {{ .Values.config.audit.query.enabled }}
    maxResults: {{ .Values.config.audit.query.maxResults }}
    timeRangeDays: {{ .Values.config.audit.query.timeRangeDays }}
  alerts:
    enabled: {{ .Values.config.audit.alerts.enabled }}
    failureThreshold: {{ .Values.config.audit.alerts.failureThreshold }}
    checkIntervalSec: {{ .Values.config.audit.alerts.checkIntervalSec }}
    method: {{ .Values.config.audit.alerts.method | quote }}
    webhookURL: {{ .Values.config.audit.alerts.webhookURL | quote }}
{{- end }}
