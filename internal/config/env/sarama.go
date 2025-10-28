package env

import (
	"github.com/IBM/sarama"
	"github.com/caarlos0/env/v11"
)

type saramaEnvConfig struct {
	GroupID string `env:"USER_CONSUMER_GROUP_ID,required"`
}

type saramaConfig struct {
	raw saramaEnvConfig
}

func NewSaramaConfig() (*saramaConfig, error) {
	var raw saramaEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}
	return &saramaConfig{raw: raw}, nil
}

func (cfg *saramaConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V4_0_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}

func (cfg *saramaConfig) GroupID() string {
	return cfg.raw.GroupID
}
