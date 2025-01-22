FROM golang:alpine AS builder

# 安装必要的工具和构建依赖
RUN apk add --no-cache curl git build-base

# 设置工作目录为项目根目录
WORKDIR /app/grpc-demo

# 复制整个项目
COPY . .

# 编译
RUN make etcd

# 最终阶段
FROM alpine:latest

WORKDIR /app

# 创建配置文件目录
RUN mkdir -p /app/configs

# 设置环境变量（你可以根据环境切换配置）
ARG ENVIRONMENT=dev

# 从构建阶段复制编译好的二进制文件和配置文件到运行阶段
COPY --from=builder /app/grpc-demo/bin/etcd-config-init .
COPY --from=builder /app/grpc-demo/configs/user-service/config.$ENVIRONMENT.yaml /app/configs/user-service/
COPY --from=builder /app/grpc-demo/configs/api-gateway/config.$ENVIRONMENT.yaml /app/configs/api-gateway/

CMD ["./etcd-config-init"]