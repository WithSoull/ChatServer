package config

import (
	"time"

	"github.com/IBM/sarama"
)

type GRPCConfig interface {
	Address() string
}

type HTTPConfig interface {
	Address() string
}

type PGConfig interface {
	DSN() string
	Timeout() time.Duration
	NeedLog() bool
}

type LoggerConfig interface {
	LogLevel() string
	AsJSON() bool
	EnableOLTP() bool
	ServiceName() string
	OTLPEndpoint() string
	ServiceEnvironment() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}

type MetricsConfig interface {
	ServiceName() string
	ServiceVersion() string
	OTLPEndpoint() string
	ServiceEnvironment() string
	PushTimeout() time.Duration
}

type RateLimiterConfig interface {
	Limit() int64
	Period() time.Duration
}

type StreamingConfig interface {
	BufferSize() int64
}

type SaramaConfig interface {
	Config() *sarama.Config
	GroupID() string
}

type KafkaConfig interface {
	Brokers() []string
}

type UserCreatedConsumerConfig interface {
	Topic() string
}

type UserDeletedConsumerConfig interface {
	Topic() string
}

type RedisConfig interface {
	Address() string
	MaxIdle() int8
	ConnTimeout() time.Duration
	IdleTimeout() time.Duration
}
