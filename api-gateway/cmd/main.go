package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/teakingwang/grpc-demo/pkg/discovery"
	userpb "github.com/teakingwang/grpc-demo/proto/user"
)

// 添加一个认证中间件
func authMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 从请求头获取 Authorization token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// 如果没有 token，设置一个测试用的 token
			r.Header.Set("Authorization", "Bearer test-token")
		} else if !strings.HasPrefix(authHeader, "Bearer ") {
			// 确保 token 格式正确
			r.Header.Set("Authorization", "Bearer "+authHeader)
		}

		// 继续处理请求
		h.ServeHTTP(w, r)
	})
}

func main() {
	// 创建根上下文
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 创建 gRPC-Gateway 的 mux
	mux := runtime.NewServeMux()

	// 创建服务发现客户端
	disc, err := discovery.NewServiceDiscovery([]string{"etcd:2379"})
	if err != nil {
		log.Printf("Failed to create service discovery: %v", err)
		log.Fatal(err)
	}
	log.Printf("Service discovery created successfully")

	/** 以下代码是 API Gateway 的主要逻辑 **/
	// 发现用户服务实例
	log.Printf("Attempting to discover user service...")
	userServices, err := disc.GetService(ctx, "user-service")
	if err != nil {
		log.Printf("Error getting user service: %v", err)
		log.Fatal("user service not found")
	}
	if len(userServices) == 0 {
		log.Printf("No user service instances found in discovery")
		log.Fatal("no user service instances found")
	}
	// 构造用户服务地址
	userAddr := fmt.Sprintf("%s:%d", userServices[0].Address, userServices[0].Port)
	log.Printf("User service discovered at address: %s", userAddr)

	// 注册用户服务的 HTTP 处理器
	err = userpb.RegisterUserServiceHandlerFromEndpoint(
		ctx, mux, userAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())},
	)
	if err != nil {
		log.Fatal(err)
	}

	// 其他微服务参照上面的代码添加即可

	// 启动 HTTP 服务器
	log.Printf("API Gateway starting on :8080")
	handler := authMiddleware(mux)
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
