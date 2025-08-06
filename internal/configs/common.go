package configs

type KafkaConfig struct {
	Brokers []string `yaml:"brokers"`
	Topic   string   `yaml:"topic"`
}

type TelemetryCfg struct {
	JaegerURL string `yaml:"jaegerURL"`
}
