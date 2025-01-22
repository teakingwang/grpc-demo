# 构建阶段
FROM golang:alpine AS builder

# 安装必要的工具和构建依赖
RUN apk add --no-cache curl git build-base

# 下载 grpc_health_probe
RUN GRPC_HEALTH_PROBE_VERSION=v0.4.19 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# 设置工作目录为项目根目录
WORKDIR /app/grpc-demo

# 复制整个项目
COPY . .

# 编译
RUN make message-service

# 最终阶段
FROM alpine:latest

WORKDIR /app

# 从 builder 阶段复制编译好的二进制文件和健康检查工具
COPY --from=builder /app/grpc-demo/bin/message-service .
COPY --from=builder /bin/grpc_health_probe /bin/grpc_health_probe

# 创建配置文件目录
RUN mkdir -p /app/configs

# 从构建阶段复制编译好的二进制文件和配置文件到运行阶段
COPY --from=builder /app/grpc-demo/bin/message-service .

EXPOSE 50052

CMD ["./message-service"] 