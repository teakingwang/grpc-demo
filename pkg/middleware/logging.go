package middleware

import (
	"context"
	"time"

	"github.com/teakingwang/grpc-demo/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// LoggingInterceptor 日志中间件,用于记录每个gRPC请求的详细信息
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 记录请求开始时间
	start := time.Now()

	// 记录请求信息
	logger.Info("Received request",
		zap.String("method", info.FullMethod), // gRPC方法名
		zap.Any("request", req),               // 请求参数
	)

	// 调用下一个处理器
	resp, err := handler(ctx, req)

	// 记录响应信息
	logger.Info("Sending response",
		zap.String("method", info.FullMethod),                                // gRPC方法名
		zap.Float64("latency_ms", float64(time.Since(start).Milliseconds())), // 请求耗时
		zap.Any("response", resp),                                            // 响应内容
		zap.Error(err),                                                       // 错误信息(如果有)
	)

	return resp, err
}
