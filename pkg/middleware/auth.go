package middleware

import (
	"context"
	"strings"

	"github.com/teakingwang/grpc-demo/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor 认证中间件,用于验证请求的认证信息
func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 跳过健康检查的认证
	if info.FullMethod == "/grpc.health.v1.Health/Check" {
		return handler(ctx, req)
	}

	// 从context中获取metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.ErrUnauthenticated
	}

	// 获取认证token
	authorization := md.Get("authorization")
	if len(authorization) == 0 {
		return nil, errors.ErrUnauthenticated
	}

	// 解析Bearer token
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	if token == "" {
		return nil, errors.ErrUnauthenticated
	}

	// TODO: 验证token
	// 这里应该添加实际的token验证逻辑,比如:
	// - 验证token格式
	// - 检查token是否过期
	// - 验证token签名
	// - 检查token权限

	return handler(ctx, req)
}
