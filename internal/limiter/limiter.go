package limiter

import (
	"context"
	"fmt"

	"rate-limiter/internal/storage"
)

type Config struct {
	RequestsPerSecond int
	BlockDuration     int
}

type RateLimiter struct {
	storage storage.Storage
	config  Config
}

func NewRateLimiter(storage storage.Storage, config Config) *RateLimiter {
	return &RateLimiter{
		storage: storage,
		config:  config,
	}
}

func (rl *RateLimiter) CheckLimit(ctx context.Context, key string) (bool, error) {
	// Incrementa o contador
	count, err := rl.storage.Increment(ctx, key)
	if err != nil {
		return false, fmt.Errorf("erro ao incrementar contador: %w", err)
	}

	// Se for a primeira requisição, define a expiração
	if count == 1 {
		err = rl.storage.Set(ctx, key, 1, rl.config.BlockDuration)
		if err != nil {
			return false, fmt.Errorf("erro ao definir expiração: %w", err)
		}
	}

	// Verifica se excedeu o limite
	if count > int64(rl.config.RequestsPerSecond) {
		return false, nil
	}

	return true, nil
}

func (rl *RateLimiter) IsBlocked(ctx context.Context, key string) (bool, error) {
	count, err := rl.storage.Get(ctx, key)
	if err != nil {
		return false, fmt.Errorf("erro ao verificar bloqueio: %w", err)
	}

	return count > int64(rl.config.RequestsPerSecond), nil
}
