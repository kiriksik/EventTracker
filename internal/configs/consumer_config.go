package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ConsumerConfig struct {
	ServiceName string       `yaml:"serviceName"`
	Kafka       KafkaConfig  `yaml:"kafka"`
	Telemetry   TelemetryCfg `yaml:"telemetry"`
}

func LoadConsumerConfig(path string) (*ConsumerConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg ConsumerConfig
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
