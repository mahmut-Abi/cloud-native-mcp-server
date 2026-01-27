FROM golang:1.25-alpine AS builder
ARG BUILD_DATE
ARG VCS_REF
ARG VERSION
RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s -X 'main.version=${VERSION}' -X 'main.commit=${VCS_REF}' -X 'main.buildTime=${BUILD_DATE}'" \
    -o k8s-mcp-server \
    ./cmd/server

FROM alpine:3.21

# Upgrade all packages to fix CVE-2024-58251 and other security vulnerabilities
RUN apk upgrade --no-cache

COPY --from=builder /app/k8s-mcp-server /usr/local/bin/k8s-mcp-server

# Create directory for audit logs with proper permissions
RUN mkdir -p /var/log/k8s-mcp-server && \
    chown 1000:1000 /var/log/k8s-mcp-server

USER root
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/k8s-mcp-server"]
