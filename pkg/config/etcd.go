package config

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Cfg struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
	Etcd     Etcd     `yaml:"etcd"`
	Service  Service  `yaml:"service"`
	Metrics  Metrics  `yaml:"metrics"`
	Log      Log      `yaml:"log"`
}

type Etcd struct {
	RegisterInterval string   `yaml:"registerInterval"`
	Endpoints        []string `yaml:"endpoints"`
	DialTimeout      string   `yaml:"dialTimeout"`
	RegisterTTL      string   `yaml:"registerTTL"`
}

type Service struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
}

type Metrics struct {
	Addr string `yaml:"addr"`
}

type Log struct {
	Output string `yaml:"output"`
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

type Server struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	Timeout string `yaml:"timeout"`
}

type Database struct {
	ConnMaxLifetime string `yaml:"connMaxLifetime"`
	Port            int    `yaml:"port"`
	MaxOpenConns    int    `yaml:"maxOpenConns"`
	MaxIdleConns    int    `yaml:"maxIdleConns"`
	Password        string `yaml:"password"`
	Dbname          string `yaml:"dbname"`
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	User            string `yaml:"user"`
}

// LoadConfigFromEtcd 从 etcd 加载配置
func (c *Config) LoadConfigFromEtcd(serviceName string) (*Cfg, error) {
	// 获取环境变量
	env := os.Getenv("ENV")
	if env == "" {
		env = "dev"
	}
	log.Printf("当前环境: %s", env)

	// 从 etcd 获取配置
	ctx := context.Background()
	key := fmt.Sprintf("/config/%s/%s", serviceName, env)
	log.Printf("获取配置键: %s", key)

	resp, err := c.client.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("获取配置失败: %v", err)
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("配置不存在: %s", key)
	}

	// 解析配置
	cfg := &Cfg{}
	if err := yaml.Unmarshal(resp.Kvs[0].Value, cfg); err != nil {
		return nil, fmt.Errorf("解析配置失败: %v", err)
	}

	return cfg, nil
}
