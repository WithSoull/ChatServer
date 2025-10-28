package env

import (
	"github.com/caarlos0/env/v11"
)

type userCreatedConsumerEnvConfig struct {
	TopicName string `env:"USER_CREATED_TOPIC_NAME,required"`
}

type userCreatedConsumerConfig struct {
	raw userCreatedConsumerEnvConfig
}

func NewUserCreatedConsumerConfig() (*userCreatedConsumerConfig, error) {
	var raw userCreatedConsumerEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &userCreatedConsumerConfig{raw: raw}, nil
}

func (cfg *userCreatedConsumerConfig) Topic() string {
	return cfg.raw.TopicName
}
