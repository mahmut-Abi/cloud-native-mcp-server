---
title: "Docker 部署"
weight: 20
---

# Docker 部署

本指南描述如何使用 Docker 部署 Cloud Native MCP Server。

## 前提条件

- Docker 20.10+ 已安装
- 适当的 Docker 权限

## 快速启动

### 基本运行

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  mahmutabi/cloud-native-mcp-server:latest
```

### 自定义配置

```bash
docker run -d \
  --name cloud-native-mcp-server \
  -p 8080:8080 \
  -v ~/.kube:/root/.kube:ro \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  -e MCP_MODE=sse \
  -e MCP_ADDR=0.0.0.0:8080 \
  -e MCP_LOG_LEVEL=info \
  mahmutabi/cloud-native-mcp-server:latest
```

## Docker Compose

### 创建 docker-compose.yml

```yaml
version: '3.8'

services:
  cloud-native-mcp-server:
    image: mahmutabi/cloud-native-mcp-server:latest
    container_name: cloud-native-mcp-server
    ports:
      - "8080:8080"
    volumes:
      - ~/.kube:/root/.kube:ro
      - ./config.yaml:/app/config.yaml:ro
    environment:
      - MCP_MODE=sse
      - MCP_ADDR=0.0.0.0:8080
      - MCP_LOG_LEVEL=info
    restart: unless-stopped
    healthcheck:
      test: ["CMD", "wget", "--spider", "http://localhost:8080/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    networks:
      - monitoring

networks:
  monitoring:
    external: true
```

### 使用 Docker Compose

```bash
# 启动
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止
docker-compose down

# 重启
docker-compose restart
```

## 配置文件

### config.yaml

```yaml
server:
  mode: "sse"
  addr: "0.0.0.0:8080"

logging:
  level: "info"
  json: false

kubernetes:
  kubeconfig: ""

grafana:
  enabled: true
  url: "http://grafana:3000"
  apiKey: "${GRAFANA_API_KEY}"

prometheus:
  enabled: true
  address: "http://prometheus:9090"

auth:
  enabled: true
  mode: "apikey"
  apiKey: "${MCP_AUTH_API_KEY}"
```

## 环境变量

### 支持的环境变量

```bash
MCP_MODE=sse
MCP_ADDR=0.0.0.0:8080
MCP_LOG_LEVEL=info
MCP_AUTH_ENABLED=true
MCP_AUTH_MODE=apikey
MCP_AUTH_API_KEY=your-key
```

### 使用 .env 文件

创建 `.env` 文件：

```bash
MCP_MODE=sse
MCP_ADDR=0.0.0.0:8080
MCP_LOG_LEVEL=info
MCP_AUTH_API_KEY=your-secret-key
GRAFANA_API_KEY=your-grafana-key
```

更新 docker-compose.yml：

```yaml
services:
  cloud-native-mcp-server:
    # ... 其他配置
    env_file:
      - .env
```

## 卷挂载

### 挂载 kubeconfig

```bash
-v ~/.kube:/root/.kube:ro
```

### 挂载配置文件

```bash
-v $(pwd)/config.yaml:/app/config.yaml:ro
```

### 挂载日志目录

```bash
-v $(pwd)/logs:/var/log/cloud-native-mcp-server
```

## 网络配置

### 使用自定义网络

```yaml
networks:
  mcp-network:
    driver: bridge
```

### 连接到外部网络

```yaml
networks:
  - external_network
  - mcp-network

networks:
  external_network:
    external: true
  mcp-network:
    driver: bridge
```

## 资源限制

```yaml
services:
  cloud-native-mcp-server:
    # ... 其他配置
    deploy:
      resources:
        limits:
          cpus: '1.0'
          memory: 512M
        reservations:
          cpus: '0.5'
          memory: 256M
```

## 健康检查

### 内置健康检查

```yaml
healthcheck:
  test: ["CMD", "wget", "--spider", "http://localhost:8080/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

### 自定义健康检查

```yaml
healthcheck:
  test: ["CMD-SHELL", "curl -f http://localhost:8080/health || exit 1"]
  interval: 30s
  timeout: 10s
  retries: 3
```

## 自定义镜像

### Dockerfile

```dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o cloud-native-mcp-server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates wget
WORKDIR /root/

COPY --from=builder /app/cloud-native-mcp-server .

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
  CMD wget --spider -q http://localhost:8080/health || exit 1

CMD ["./cloud-native-mcp-server"]
```

### 构建镜像

```bash
docker build -t your-registry/cloud-native-mcp-server:latest .
```

### 推送镜像

```bash
docker push your-registry/cloud-native-mcp-server:latest
```

## 多容器部署

### 完整栈示例

```yaml
version: '3.8'

services:
  cloud-native-mcp-server:
    image: mahmutabi/cloud-native-mcp-server:latest
    ports:
      - "8080:8080"
    volumes:
      - ~/.kube:/root/.kube:ro
      - ./config.yaml:/app/config.yaml:ro
    env_file:
      - .env
    depends_on:
      - grafana
      - prometheus
    networks:
      - monitoring

  grafana:
    image: grafana/grafana:latest
    ports:
      - "3000:3000"
    volumes:
      - grafana-data:/var/lib/grafana
    networks:
      - monitoring

  prometheus:
    image: prom/prometheus:latest
    ports:
      - "9090:9090"
    volumes:
      - prometheus-data:/prometheus
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
    networks:
      - monitoring

volumes:
  grafana-data:
  prometheus-data:

networks:
  monitoring:
    driver: bridge
```

## 日志管理

### 查看日志

```bash
# 实时日志
docker-compose logs -f

# 最近 100 行
docker-compose logs --tail=100

# 特定服务
docker-compose logs -f cloud-native-mcp-server
```

### 日志驱动

```yaml
services:
  cloud-native-mcp-server:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

## 故障排查

### 检查容器状态

```bash
docker ps -a
docker logs cloud-native-mcp-server
docker inspect cloud-native-mcp-server
```

### 进入容器

```bash
docker exec -it cloud-native-mcp-server sh
```

### 重新构建

```bash
docker-compose down
docker-compose build
docker-compose up -d
```

## 相关文档

- [Kubernetes 部署](/zh/guides/deployment/kubernetes/)
- [Helm 部署](/zh/guides/deployment/helm/)
- [配置指南](/zh/guides/configuration/)
- [Docker 文档](https://docs.docker.com/)