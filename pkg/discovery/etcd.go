package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Service 表示一个微服务实例的信息
type Service struct {
	Name    string `json:"name"`    // 服务名称
	ID      string `json:"id"`      // 服务实例唯一标识
	Address string `json:"address"` // 服务地址
	Port    int    `json:"port"`    // 服务端口
}

// ServiceRegistry 服务注册器,负责将服务注册到etcd
type ServiceRegistry struct {
	client        *clientv3.Client                        // etcd客户端
	leaseID       clientv3.LeaseID                        // 租约ID
	keepAliveChan <-chan *clientv3.LeaseKeepAliveResponse // 续租通道
	key           string                                  // 服务注册的key
	service       *Service                                // 服务信息
}

// NewServiceRegistry 创建一个新的服务注册器
func NewServiceRegistry(endpoints []string, service *Service) (*ServiceRegistry, error) {
	// 创建etcd客户端
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	// 构造服务注册key
	srv := &ServiceRegistry{
		client:  client,
		key:     fmt.Sprintf("/services/%s", service.Name),
		service: service,
	}

	return srv, nil
}

// Register 注册服务到etcd
func (s *ServiceRegistry) Register() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建租约(TTL为10秒)
	leaseResp, err := s.client.Grant(ctx, 10)
	if err != nil {
		return err
	}

	s.leaseID = leaseResp.ID

	// 将服务信息序列化并注册到etcd
	serviceValue, _ := json.Marshal(s.service)
	_, err = s.client.Put(ctx, s.key, string(serviceValue), clientv3.WithLease(s.leaseID))
	if err != nil {
		return err
	}

	// 开启自动续租
	keepAliveChan, err := s.client.KeepAlive(context.Background(), s.leaseID)
	if err != nil {
		return err
	}
	s.keepAliveChan = keepAliveChan

	return nil
}

// Unregister 注销服务
func (s *ServiceRegistry) Unregister() error {
	// 撤销租约会自动删除关联的key
	if s.leaseID != 0 {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err := s.client.Revoke(ctx, s.leaseID)
		return err
	}
	return nil
}
