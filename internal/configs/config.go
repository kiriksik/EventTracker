package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ServiceName string        `yaml:"serviceName"`
	HTTP        HTTPConfig    `yaml:"http"`
	Kafka       KafkaConfig   `yaml:"kafka"`
	ClickHouse  ClickHouseCfg `yaml:"clickhouse"`
	Telemetry   TelemetryCfg  `yaml:"telemetry"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

type ClickHouseCfg struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

type TelemetryCfg struct {
	JaegerURL string `yaml:"jaegerURL"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)

	var cfg Config
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
