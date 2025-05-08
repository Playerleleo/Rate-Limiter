package storage

import "context"

// Storage define a interface para diferentes implementações de armazenamento
type Storage interface {
	// Increment incrementa o contador para uma chave específica
	Increment(ctx context.Context, key string) (int64, error)

	// Get retorna o valor atual do contador para uma chave
	Get(ctx context.Context, key string) (int64, error)

	// Set define um valor para uma chave com expiração
	Set(ctx context.Context, key string, value int64, expiration int) error

	// Delete remove uma chave
	Delete(ctx context.Context, key string) error

	// Close fecha a conexão com o storage
	Close() error
}
