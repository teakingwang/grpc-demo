package middleware

import (
	"context"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	grpcstatus "google.golang.org/grpc/status"
)

// 定义Prometheus指标
var (
	// 请求总数计数器
	grpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_requests_total",           // 指标名称
			Help: "Total number of gRPC requests", // 指标说明
		},
		[]string{"method", "status"}, // 标签:方法名和状态
	)

	// 请求耗时直方图
	grpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_duration_seconds",    // 指标名称
			Help:    "gRPC request duration in seconds", // 指标说明
			Buckets: prometheus.DefBuckets,              // 默认的bucket区间
		},
		[]string{"method"}, // 标签:方法名
	)

	// 活跃请求数量
	grpcRequestsInProgress = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "grpc_requests_in_progress",                     // 指标名称
			Help: "Number of gRPC requests currently in progress", // 指标说明
		},
		[]string{"method"}, // 标签:方法名
	)

	// 请求大小直方图
	grpcRequestSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_request_size_bytes",                          // 指标名称
			Help:    "Size of gRPC requests in bytes",                   // 指标说明
			Buckets: []float64{32, 64, 128, 256, 512, 1024, 2048, 4096}, // 自定义bucket区间
		},
		[]string{"method"}, // 标签:方法名
	)

	// 响应大小直方图
	grpcResponseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_response_size_bytes",                         // 指标名称
			Help:    "Size of gRPC responses in bytes",                  // 指标说明
			Buckets: []float64{32, 64, 128, 256, 512, 1024, 2048, 4096}, // 自定义bucket区间
		},
		[]string{"method"}, // 标签:方法名
	)
)

// 初始化:注册所有指标到Prometheus
func init() {
	prometheus.MustRegister(grpcRequestsTotal)
	prometheus.MustRegister(grpcRequestDuration)
	prometheus.MustRegister(grpcRequestsInProgress)
	prometheus.MustRegister(grpcRequestSize)
	prometheus.MustRegister(grpcResponseSize)
}

// MetricsInterceptor 指标收集中间件,用于收集gRPC请求的各项指标
func MetricsInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	// 增加活跃请求计数
	grpcRequestsInProgress.WithLabelValues(info.FullMethod).Inc()
	defer grpcRequestsInProgress.WithLabelValues(info.FullMethod).Dec()

	// 记录请求大小
	if msg, ok := req.(proto.Message); ok {
		grpcRequestSize.WithLabelValues(info.FullMethod).Observe(float64(proto.Size(msg)))
	}

	// 处理请求
	resp, err := handler(ctx, req)

	// 记录请求耗时
	duration := time.Since(start).Seconds()
	grpcRequestDuration.WithLabelValues(info.FullMethod).Observe(duration)

	// 记录响应大小
	if msg, ok := resp.(proto.Message); ok {
		grpcResponseSize.WithLabelValues(info.FullMethod).Observe(float64(proto.Size(msg)))
	}

	// 记录请求结果
	result := "success"
	if err != nil {
		if st, ok := grpcstatus.FromError(err); ok {
			result = st.Code().String()
		} else {
			result = "error"
		}
	}
	grpcRequestsTotal.WithLabelValues(info.FullMethod, result).Inc()

	return resp, err
}
