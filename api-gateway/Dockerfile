# 第一阶段：构建阶段，使用 golang:alpine 作为基础镜像
FROM golang:alpine AS builder

# 设置构建阶段的工作目录为 /app
WORKDIR /app/grpc-demo

# 安装构建过程中需要的工具：make 用于执行 Makefile，git 用于版本信息
RUN apk add --no-cache make git

# 将当前目录下的所有文件复制到容器的工作目录中
COPY . .

# 下载项目所需的 Go 依赖包
RUN go mod download
# 使用 Makefile 中定义的命令构建 api-gateway 服务
RUN make api-gateway

# 第二阶段：运行阶段，使用轻量级的 alpine 作为基础镜像
FROM alpine:latest
RUN apk add --no-cache netcat-openbsd curl busybox-extras
WORKDIR /app

# 创建证书目录
RUN mkdir -p /app/certs

# 从构建阶段复制编译好的二进制文件和配置文件到运行阶段
COPY --from=builder /app/grpc-demo/bin/api-gateway .
COPY --from=builder /app/grpc-demo/certs/* ./certs/

# 声明容器将监听的端口号
EXPOSE 8080

# 设置容器启动时执行的命令
ENTRYPOINT ["./api-gateway"] 