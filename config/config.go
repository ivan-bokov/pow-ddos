package config

import (
	"context"
	"fmt"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type ServerConfig struct {
	Addr         string        `env:"SERVER_ADDR,default=0.0.0.0:8000"`
	KeepAlive    time.Duration `env:"SERVER_KEEPALIVE,default=20s"`
	WriteTimeout time.Duration `env:"SERVER_WRITE_TIMEOUT,default=100s"`
	ReadTimeout  time.Duration `env:"SERVER_READ_TIMEOUT,default=100s"`
	IdleTimeout  time.Duration `env:"SERVER_IDLE_TIMEOUT,default=100s"`
	Pow          PowConfig     `env:",prefix=POW_"`
	Logger       LoggerConfig  `env:",prefix=LOGGER_"`
	DDOS         DdosConfig    `env:",prefix=DDOS_"`
	Storage      StorageConfig `env:",prefix=STORAGE_"`
}

type LoggerConfig struct {
	Level string `env:"LEVEL,default=info"`
}

type PowConfig struct {
	ChallengeLength int   `env:"CHALLENGE_LENGTH,default=8"`
	DifficultyStart int   `env:"DIFFICULTY,default=3"`
	MaxDifficulty   uint8 `env:"MAX_DIFFICULTY,default=10"`
}
type DdosConfig struct {
	Window time.Duration `env:"WINDOW,default=2s"`
	Rate   uint64        `env:"RATE,default=10"`
}
type StorageConfig struct {
	Path string `env:"PATH,default=data/word-of-wisdom.txt"`
}

func NewConfig[T any](ctx context.Context) (T, error) {
	var config T
	if err := envconfig.Process(ctx, &config); err != nil {
		return config, fmt.Errorf("failed to process config: %w", err)
	}

	return config, nil
}
