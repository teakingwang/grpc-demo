.PHONY: all proto clean user-service api-gateway build-all test lint deps etcd

# 设置 Go 编译器参数
GOPATH:=$(shell go env GOPATH)
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

# 设置输出目录
OUT_DIR=bin
PROTO_DIR=proto

# 设置服务名称和输出路径
USER_SERVICE=user-service
API_GATEWAY=api-gateway
ETCD_CONFIG_INIT=etcd-config-init
USER_SERVICE_PATH=user-service/cmd
API_GATEWAY_PATH=api-gateway/cmd
ETCD_CONFIG_INIT_PATH=tools/etcd-config-init

# 自动检测操作系统
GOOS := $(shell go env GOOS)
GOARCH := $(shell go env GOARCH)

# 添加版本信息
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME)

all: proto user-service api-gateway etcd

# 生成 proto 文件， --proto_path可以使用-I代替
proto:
	@echo "生成 proto 文件..."
	@mkdir -p $(PROTO_DIR)/user/gen
	protoc --proto_path=$(PROTO_DIR)/user \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTO_DIR)/user/gen --go_opt=paths=source_relative \
		--go-grpc_out=$(PROTO_DIR)/user/gen --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(PROTO_DIR)/user/gen --grpc-gateway_opt=paths=source_relative \
		$(PROTO_DIR)/user/*.proto
# 如果有其他微服务proto，安装相同的方式生成即可

# 构建 user-service
user-service:
	@echo "构建 user-service..."
	@mkdir -p $(OUT_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)"  -o $(OUT_DIR)/$(USER_SERVICE) ./$(USER_SERVICE_PATH)

# 构建 api-gateway
api-gateway:
	@echo "构建 api-gateway..."
	@mkdir -p $(OUT_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(OUT_DIR)/$(API_GATEWAY) ./$(API_GATEWAY_PATH)

# etcd
etcd:
	@echo "构建 etcd..."
	@mkdir -p $(OUT_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(OUT_DIR)/$(ETCD_CONFIG_INIT) ./$(ETCD_CONFIG_INIT_PATH)


# 清理构建文件
clean:
	@echo "清理构建文件..."
	@rm -rf $(OUT_DIR)
	$(GOCLEAN) 

build-all: user-service api-gateway
	# 这里可以添加其他服务的构建命令 

# 添加测试命令
test:
	go test -v ./...

# 添加代码检查
lint:
	golangci-lint run

# 添加依赖安装
deps:
	go mod download
	go mod tidy 