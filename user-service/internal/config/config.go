package config

type Config struct {
	Server struct {
		Port int    `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`

	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
	} `yaml:"database"`

	Etcd struct {
		ServiceName string `yaml:"serviceName"`
		NodeName    string `yaml:"nodeName"`
		TTL         int64  `yaml:"ttl"`
	} `yaml:"etcd"`
}
