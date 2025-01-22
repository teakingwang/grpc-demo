package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/teakingwang/grpc-demo/message-service/internal/service"
	"github.com/teakingwang/grpc-demo/pkg/config"
	"github.com/teakingwang/grpc-demo/pkg/discovery"
	"github.com/teakingwang/grpc-demo/pkg/health"
	"github.com/teakingwang/grpc-demo/pkg/logger"
	"github.com/teakingwang/grpc-demo/pkg/middleware"
	pb "github.com/teakingwang/grpc-demo/proto/message"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// 与 user-service/cmd/main.go 类似,只是端口和服务名不同
// 已经在前面添加过注释,这里不再重复

func main() {
	// todo换成读取配置文件
	etcdEndpoint := "http://etcd:2379"

	// 初始化配置管理器,用于从etcd获取和保存配置
	cfgCli, err := config.NewConfig([]string{etcdEndpoint})
	if err != nil {
		logger.Fatal("Failed to create config", zap.Error(err))
	}

	cfg, err := cfgCli.LoadConfigFromEtcd("message-service")
	if err != nil {
		logger.Fatal("Fail to LoadConfigFromEtcd", zap.Error(err))
	}

	// 连接MySQL数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Dbname)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// 创建gRPC服务器并注册中间件
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor, // 日志中间件
			middleware.AuthInterceptor,    // 认证中间件
			middleware.MetricsInterceptor, // 指标收集中间件
		),
	)

	// 注册消息服务
	messageService := service.NewMessageService(db)
	pb.RegisterMessageServiceServer(grpcServer, messageService)

	// 注册健康检查服务
	healthChecker := health.NewHealthChecker()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthChecker)

	// 启动gRPC服务器
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// 注册服务到etcd
	registry, err := discovery.NewServiceRegistry([]string{etcdEndpoint}, &discovery.Service{
		Name:    "message-service",
		ID:      "message-1",
		Address: "message-service",
		Port:    50052,
	})
	if err != nil {
		log.Fatalf("Failed to create service registry: %v", err)
	}

	// 注册服务
	if err := registry.Register(); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	defer registry.Unregister()

	// 启动 Prometheus 指标服务器
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":9090", mux); err != nil {
			logger.Fatal("Failed to start metrics server", zap.Error(err))
		}
	}()

	// 优雅关闭处理
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
		grpcServer.GracefulStop()
	}()

	// 启动服务
	log.Printf("Message service starting on :50052")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
