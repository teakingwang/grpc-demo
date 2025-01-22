package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type ServiceDiscovery struct {
	client *clientv3.Client
}

const (
	maxRetries = 5
	retryDelay = time.Second * 3
)

func NewServiceDiscovery(endpoints []string) (*ServiceDiscovery, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("创建 etcd 客户端失败: %v", err)
	}
	return &ServiceDiscovery{client: client}, nil
}

func (sd *ServiceDiscovery) GetService(ctx context.Context, serviceName string) ([]ServiceInstance, error) {
	key := fmt.Sprintf("/services/%s", serviceName)

	for i := 0; i < maxRetries; i++ {
		log.Printf("[尝试 %d/%d] 查找服务: %s", i+1, maxRetries, key)

		// 列出所有服务
		resp, err := sd.client.Get(ctx, "/services/", clientv3.WithPrefix())
		if err != nil {
			log.Printf("获取服务列表失败: %v", err)
		} else {
			log.Printf("当前注册的所有服务:")
			for _, kv := range resp.Kvs {
				log.Printf("  %s -> %s", string(kv.Key), string(kv.Value))
			}
		}

		// 获取特定服务
		resp, err = sd.client.Get(ctx, key)
		if err != nil {
			log.Printf("获取服务失败: %v", err)
			time.Sleep(retryDelay)
			continue
		}

		if len(resp.Kvs) == 0 {
			log.Printf("服务未注册，等待 %v 后重试...", retryDelay)
			time.Sleep(retryDelay)
			continue
		}

		var instances []ServiceInstance
		for _, kv := range resp.Kvs {
			var instance ServiceInstance
			if err := json.Unmarshal(kv.Value, &instance); err != nil {
				log.Printf("解析服务信息失败: %v, value: %s", err, string(kv.Value))
				continue
			}
			instances = append(instances, instance)
			log.Printf("发现服务实例: %+v", instance)
		}

		if len(instances) > 0 {
			return instances, nil
		}
	}

	return nil, fmt.Errorf("服务发现失败: 未找到服务 %s", serviceName)
}

type ServiceInstance struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// RegisterService 注册服务到 etcd
func RegisterService(serviceName, serviceAddr string, port int, ttl int64) error {
	log.Printf("开始注册服务: name=%s, addr=%s, port=%d", serviceName, serviceAddr, port)

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"etcd:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return fmt.Errorf("连接 etcd 失败: %v", err)
	}
	defer cli.Close()

	// 服务信息
	serviceInfo := ServiceInstance{
		Name:    serviceName,
		Address: serviceAddr,
		Port:    port,
	}

	key := fmt.Sprintf("/services/%s", serviceName)
	value, _ := json.Marshal(serviceInfo)

	// 先删除旧的注册信息
	ctx := context.Background()
	_, err = cli.Delete(ctx, key)
	if err != nil {
		log.Printf("Warning: failed to delete old key: %v", err)
	}

	// 注册新的服务信息
	_, err = cli.Put(ctx, key, string(value))
	if err != nil {
		return fmt.Errorf("注册服务失败: %v", err)
	}

	// 验证注册是否成功
	resp, err := cli.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("验证注册失败: %v", err)
	}
	if len(resp.Kvs) == 0 {
		return fmt.Errorf("服务注册后未找到: %s", key)
	}

	log.Printf("服务注册成功: %s -> %s", key, string(value))
	return nil
}
