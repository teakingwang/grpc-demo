package config

import (
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Config 配置管理器
type Config struct {
	client *clientv3.Client
}

// NewConfig 创建新的配置管理器实例
func NewConfig(endpoints []string) (*Config, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})

	// 如果连接etcd失败，仅记录错误但不返回错误
	if err != nil {
		fmt.Printf("Failed to connect to etcd: %v, will use local config files\n", err)
		client = nil
	}

	return &Config{
		client: client,
	}, nil
}
