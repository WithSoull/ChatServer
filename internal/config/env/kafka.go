package env

import (
	"strings"

	"github.com/caarlos0/env/v11"
)

type kafkaEnvConfig struct {
	Brokers string `env:"KAFKA_BROKERS,required"`
}

type kafkaConfig struct {
	raw kafkaEnvConfig
}

func NewKafkaConfig() (*kafkaConfig, error) {
	var raw kafkaEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &kafkaConfig{raw: raw}, nil
}

func (cfg *kafkaConfig) Brokers() []string {
	return strings.Split(cfg.raw.Brokers, ",")
}
