package env

import (
	"github.com/caarlos0/env/v11"
)

type userDeletedConsumerEnvConfig struct {
	TopicName string `env:"USER_DELETED_TOPIC_NAME,required"`
}

type userDeletedConsumerConfig struct {
	raw userDeletedConsumerEnvConfig
}

func NewUserDeletedConsumerConfig() (*userDeletedConsumerConfig, error) {
	var raw userDeletedConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &userDeletedConsumerConfig{raw: raw}, nil
}

func (cfg *userDeletedConsumerConfig) Topic() string {
	return cfg.raw.TopicName
}
