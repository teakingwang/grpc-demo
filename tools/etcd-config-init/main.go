package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
)

var (
	configDir = flag.String("config-dir", "configs", "配置文件目录路径")
	env       = flag.String("env", "dev", "环境: dev, test, prod")
)

func main() {
	flag.Parse()
	log.Println("开始初始化配置")

	// 创建 etcd 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"etcd:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// 遍历配置目录
	services := []string{"api-gateway", "user-service", "message-service"}
	for _, service := range services {
		configPath := filepath.Join(*configDir, service, fmt.Sprintf("config.%s.yaml", *env))

		// 读取配置文件
		data, err := os.ReadFile(configPath)
		if err != nil {
			log.Printf("读取配置文件失败 %s: %v", configPath, err)
			continue
		}

		// 解析 YAML
		var config map[string]interface{}
		if err := yaml.Unmarshal(data, &config); err != nil {
			log.Printf("解析 YAML 失败 %s: %v", configPath, err)
			continue
		}

		// 转换为 JSON
		jsonData, err := yaml.Marshal(config)
		if err != nil {
			log.Printf("转换 JSON 失败 %s: %v", configPath, err)
			continue
		}

		// 存储到 etcd
		key := fmt.Sprintf("/config/%s/%s", service, *env)
		ctx := context.Background()
		_, err = cli.Put(ctx, key, string(jsonData))
		if err != nil {
			log.Printf("存储到 etcd 失败 %s: %v", key, err)
			continue
		}

		log.Printf("成功初始化配置: %s", key)
	}

	printAllConfigs(cli)
}

// 打印当前 etcd 中的所有配置
func printAllConfigs(cli *clientv3.Client) {
	ctx := context.Background()
	resp, err := cli.Get(ctx, "/config/", clientv3.WithPrefix())
	if err != nil {
		log.Printf("获取配置失败: %v", err)
		return
	}

	fmt.Println("\n当前 etcd 中的配置:")
	for _, kv := range resp.Kvs {
		fmt.Printf("\nKey: %s\nValue: %s\n", kv.Key, kv.Value)
		fmt.Println(strings.Repeat("-", 50))
	}
}
