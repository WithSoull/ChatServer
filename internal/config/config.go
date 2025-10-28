package config

import (
	"os"

	"github.com/WithSoull/ChatServer/internal/config/env"
	"github.com/joho/godotenv"
)

// appConfig holds the global application configuration instance.
var appConfig *config

// config represents the complete application configuration.
type config struct {
	Logger      LoggerConfig
	GRPC        GRPCConfig
	HTTP        HTTPConfig
	PG          PGConfig
	Tracing     TracingConfig
	Metrics     MetricsConfig
	RateLimiter RateLimiterConfig
	Streaming   StreamingConfig
	Redis       RedisConfig

	Kafka               KafkaConfig
	Sarama              SaramaConfig
	UserCreatedConsumer UserCreatedConsumerConfig
	UserDeletedConsumer UserDeletedConsumerConfig
}

// Load reads environment variables from .env file(s) and initializes the application configuration.
// The function ignores file-not-found errors, allowing configuration to be loaded
// from system environment variables when .env file is absent.
func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	grpcCfg, err := env.NewGRPCConfig()
	if err != nil {
		return err
	}

	httpCfg, err := env.NewHTTPConfig()
	if err != nil {
		return err
	}

	pgCfg, err := env.NewPGConfig()
	if err != nil {
		return err
	}

	tracingCfg, err := env.NewTracingConfig()
	if err != nil {
		return err
	}

	metricsCfg, err := env.NewMetricsConfig()
	if err != nil {
		return err
	}

	rateLimiterCfg, err := env.NewRateLimiterConfig()
	if err != nil {
		return err
	}

	streamingCfg, err := env.NewStreamingConfig()
	if err != nil {
		return err
	}

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	saramaCfg, err := env.NewSaramaConfig()
	if err != nil {
		return err
	}

	userCreatedConsumerCfg, err := env.NewUserCreatedConsumerConfig()
	if err != nil {
		return err
	}

	userDeletedConsumerCfg, err := env.NewUserDeletedConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:      loggerCfg,
		GRPC:        grpcCfg,
		HTTP:        httpCfg,
		PG:          pgCfg,
		Tracing:     tracingCfg,
		Metrics:     metricsCfg,
		RateLimiter: rateLimiterCfg,
		Streaming:   streamingCfg,
		Redis:       redisCfg,

		Kafka:               kafkaCfg,
		Sarama:              saramaCfg,
		UserCreatedConsumer: userCreatedConsumerCfg,
		UserDeletedConsumer: userDeletedConsumerCfg,
	}

	return nil
}

// AppConfig returns the global application configuration instance.
// Panics if called before Load().
func AppConfig() *config {
	return appConfig
}
