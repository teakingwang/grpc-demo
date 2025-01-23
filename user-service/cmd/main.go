package main

import (
	"database/sql"
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/teakingwang/grpc-demo/pkg/config"
	"github.com/teakingwang/grpc-demo/pkg/discovery"
	"github.com/teakingwang/grpc-demo/pkg/health"
	"github.com/teakingwang/grpc-demo/pkg/logger"
	"github.com/teakingwang/grpc-demo/pkg/middleware"
	pb "github.com/teakingwang/grpc-demo/proto/user/gen"
	"github.com/teakingwang/grpc-demo/user-service/internal/service"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {
	// todo换成读取配置文件
	etcdEndpoint := "http://etcd:2379"

	// 初始化配置管理器,用于从etcd获取和保存配置
	cfgCli, err := config.NewConfig([]string{etcdEndpoint})
	if err != nil {
		logger.Fatal("Failed to create config", zap.Error(err))
	}

	cfg, err := cfgCli.LoadConfigFromEtcd("user-service")
	if err != nil {
		logger.Fatal("Fail to LoadConfigFromEtcd", zap.Error(err))
	}

	// 连接MySQL数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Dbname)
	db, err := sql.Open("mysql", dsn)
	log.Println(dsn)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// 启用 TLS
	//certFile := "./certs/server.crt"
	//keyFile := "./certs/server.key"
	//creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
	//if err != nil {
	//	log.Fatalf("failed to load server TLS certificates: %v", err)
	//}

	// 创建gRPC服务器并注册中间件
	grpcServer := grpc.NewServer(
		//grpc.Creds(creds), // 启用TLS
		grpc.ChainUnaryInterceptor(
			middleware.LoggingInterceptor, // 日志中间件
			middleware.AuthInterceptor,    // 认证中间件
			middleware.MetricsInterceptor, // 指标收集中间件
		),
	)

	// 注册用户服务
	userService := service.NewUserService(db)
	pb.RegisterUserServiceServer(grpcServer, userService)
	reflection.Register(grpcServer)

	// 注册健康检查服务
	healthChecker := health.NewHealthChecker()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthChecker)

	// 启动gRPC服务器监听
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// 创建服务注册器并注册到etcd
	registry, err := discovery.NewServiceRegistry([]string{etcdEndpoint}, &discovery.Service{
		Name:    "user-service", // 服务名称
		ID:      "user-1",       // 服务实例ID
		Address: "user-service", // 服务地址
		Port:    50051,          // 服务端口
	})
	if err != nil {
		log.Fatalf("Failed to create service registry: %v", err)
	}

	// 注册服务并确保服务退出时注销
	if err := registry.Register(); err != nil {
		log.Fatalf("Failed to register service: %v", err)
	}
	defer registry.Unregister()

	// 启动Prometheus指标收集服务器
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
		grpcServer.GracefulStop() // 优雅停止gRPC服务器
	}()

	// 启动服务并阻塞等待
	log.Printf("User service starting on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
