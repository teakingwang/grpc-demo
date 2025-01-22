package health

import (
	"context"

	"google.golang.org/grpc/health/grpc_health_v1"
)

// HealthChecker 实现gRPC健康检查服务
type HealthChecker struct{}

// NewHealthChecker 创建新的健康检查器
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{}
}

// Check 实现健康检查接口
// 返回服务的健康状态
func (h *HealthChecker) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	// 当前简单实现:始终返回健康状态
	// 实际使用时可以根据服务的具体状态返回不同的结果
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch 实现健康状态监控接口
// 用于监控服务的健康状态变化
func (h *HealthChecker) Watch(req *grpc_health_v1.HealthCheckRequest, server grpc_health_v1.Health_WatchServer) error {
	// 当前简单实现:发送一次健康状态后返回
	// 实际使用时可以持续监控服务状态并发送状态变更
	return server.Send(&grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	})
}
