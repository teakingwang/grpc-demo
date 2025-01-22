# grpc-demo

一个使用 gRPC 和 grpc-gateway 的示例项目，支持 gRPC 和 RESTful API 双协议。

## 项目结构

```tree
.
├── LICENSE
├── Makefile                 # 项目构建脚本
├── README.md                # 项目说明文档
├── docker-compose.yaml      # Docker 编排文件
├── go.mod                   # Go 模块依赖
├── proto                    # Protocol Buffers 定义目录
│   └── user                 # 用户服务相关协议
│       ├── user.proto       # 用户服务协议定义
│       └── user.pb.gw.go    # gRPC-Gateway 生成的代码
├── third_party              # 第三方依赖
│   └── googleapis          # Google API 定义
│       └── google
│           └── api
│               ├── annotations.proto
│               └── http.proto
└── alertmanager            # 告警管理配置
    └── alertmanager.yml    # Alertmanager 配置文件

## 功能特性

- 支持 gRPC 协议调用
- 通过 grpc-gateway 支持 RESTful API 调用
- 包含完整的 Protocol Buffers 定义
- 集成 Alertmanager 告警管理

## 快速开始

### 环境要求

- Go 1.16+
- Protocol Buffers 编译器
- Docker 和 Docker Compose（可选）

### 安装依赖

```bash:README.md
# 安装 protoc 编译器
brew install protobuf  # MacOS
# 或
apt-get install protobuf-compiler  # Ubuntu

# 安装 Go 相关工具
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```

### 编译项目

```bash
# 生成 Protocol Buffers 代码
make proto

# 构建项目
make build

# etcd健康检查
etcdctl --endpoints=http://etcd:2379 endpoint health

# 测试grpc
brew install grpcurl (mac安装)
grpcurl -H "Authorization: Bearer your_token_value" -d '{"id": "1"}' -plaintext localhost:50051 user.UserService.GetUser

# 测试grpcgateway
curl -H "Authorization: Bearer your_token_value" http://localhost:8080/v1/user/1
```

### 使用 Docker 运行

```bash:README.md
docker-compose up -d
```

## API 文档

### gRPC 接口

用户服务定义在 `proto/user/user.proto` 文件中，包含以下接口：

- 创建用户
- 获取用户信息
- 更新用户信息
- 删除用户

### RESTful API

通过 grpc-gateway，所有 gRPC 接口都会自动生成对应的 HTTP 接口。

## 贡献指南

欢迎提交 Issue 和 Pull Request。

## 许可证

本项目采用 [LICENSE](./LICENSE) 协议。