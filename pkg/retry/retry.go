package retry

import (
	"context"
	"time"
)

// RetryConfig 重试配置
type RetryConfig struct {
	MaxAttempts     int           // 最大重试次数
	InitialInterval time.Duration // 初始重试间隔
	MaxInterval     time.Duration // 最大重试间隔
}

// DefaultConfig 返回默认的重试配置
func DefaultConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts:     3,                      // 默认重试3次
		InitialInterval: 100 * time.Millisecond, // 初始间隔100ms
		MaxInterval:     1 * time.Second,        // 最大间隔1s
	}
}

// Do 执行带重试的操作
// ctx: 上下文,用于取消操作
// fn: 需要重试的函数
// config: 重试配置,如果为nil则使用默认配置
func Do(ctx context.Context, fn func() error, config *RetryConfig) error {
	if config == nil {
		config = DefaultConfig()
	}

	var err error
	interval := config.InitialInterval

	for attempt := 0; attempt < config.MaxAttempts; attempt++ {
		// 执行操作
		if err = fn(); err == nil {
			return nil
		}

		// 检查上下文是否已取消
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// 最后一次尝试不需要等待
		if attempt == config.MaxAttempts-1 {
			break
		}

		// 等待后重试
		time.Sleep(interval)

		// 增加重试间隔,但不超过最大间隔
		interval *= 2
		if interval > config.MaxInterval {
			interval = config.MaxInterval
		}
	}

	return err
}
