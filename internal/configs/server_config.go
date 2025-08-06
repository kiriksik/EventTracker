package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	ServiceName string        `yaml:"serviceName"`
	HTTP        HTTPConfig    `yaml:"http"`
	Kafka       KafkaConfig   `yaml:"kafka"`
	ClickHouse  ClickHouseCfg `yaml:"clickhouse"`
	Telemetry   TelemetryCfg  `yaml:"telemetry"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type ClickHouseCfg struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func LoadServerConfig(path string) (*ServerConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg ServerConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
