package env

import (
	"github.com/caarlos0/env/v11"
)

type streamingEnvConfig struct {
	BufferSize int64 `env:"STREAMING_BUFFER_SIZE" envDefault:"100"`
}

type streamingConfig struct {
	raw streamingEnvConfig
}

func NewStreamingConfig() (*streamingConfig, error) {
	var raw streamingEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &streamingConfig{raw: raw}, nil
}

func (cfg *streamingConfig) BufferSize() int64 {
	return cfg.raw.BufferSize
}
