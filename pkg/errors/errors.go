package errors

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 预定义的gRPC错误
var (
	// ErrNotFound 资源未找到错误
	ErrNotFound = status.Error(codes.NotFound, "resource not found")

	// ErrInvalidInput 无效输入参数错误
	ErrInvalidInput = status.Error(codes.InvalidArgument, "invalid input")

	// ErrInternal 内部服务器错误
	ErrInternal = status.Error(codes.Internal, "internal error")

	// ErrUnauthenticated 未认证错误
	ErrUnauthenticated = status.Error(codes.Unauthenticated, "unauthenticated")

	// ErrPermissionDenied 权限不足错误
	ErrPermissionDenied = status.Error(codes.PermissionDenied, "permission denied")
)

// FromError 将普通错误转换为gRPC错误
// 如果错误已经是gRPC错误则直接返回
// 否则将其包装为Internal错误
func FromError(err error) error {
	if err == nil {
		return nil
	}

	// 检查是否已经是gRPC错误
	_, ok := status.FromError(err)
	if ok {
		return err
	}

	errMsg := err.Error()

	// 对网络连接错误进行更详细的处理
	switch {
	case strings.Contains(errMsg, "connection refused"):
		return status.Error(codes.Unavailable, fmt.Sprintf("目标服务不可用或未启动: %v", err))
	case strings.Contains(errMsg, "no such host"):
		return status.Error(codes.Unavailable, fmt.Sprintf("无法解析服务地址: %v", err))
	case strings.Contains(errMsg, "context deadline exceeded"):
		return status.Error(codes.DeadlineExceeded, fmt.Sprintf("服务调用超时: %v", err))
	}

	// 包装为Internal错误
	return status.Error(codes.Internal, err.Error())
}
