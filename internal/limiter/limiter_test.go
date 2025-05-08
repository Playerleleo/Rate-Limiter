package limiter

import (
	"context"
	"testing"
)

func TestRateLimiter(t *testing.T) {
	tests := []struct {
		name           string
		requestsPerSec int
		blockDuration  int
		requests       int
		wantBlocked    bool
	}{
		{
			name:           "Dentro do limite",
			requestsPerSec: 5,
			blockDuration:  1,
			requests:       5,
			wantBlocked:    false,
		},
		{
			name:           "Excedendo o limite",
			requestsPerSec: 5,
			blockDuration:  1,
			requests:       6,
			wantBlocked:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := NewMockStorage()
			limiter := NewRateLimiter(storage, Config{
				RequestsPerSecond: tt.requestsPerSec,
				BlockDuration:     tt.blockDuration,
			})

			ctx := context.Background()
			key := "test-key"

			// Faz as requisições
			for i := 0; i < tt.requests; i++ {
				allowed, err := limiter.CheckLimit(ctx, key)
				if err != nil {
					t.Errorf("CheckLimit() error = %v", err)
					return
				}

				// Se for a última requisição, verifica se foi bloqueada
				if i == tt.requests-1 {
					if allowed == tt.wantBlocked {
						t.Errorf("CheckLimit() = %v, want %v", allowed, !tt.wantBlocked)
					}
				} else if !allowed {
					t.Errorf("CheckLimit() = %v, want true", allowed)
				}
			}

			// Verifica se está bloqueado
			blocked, err := limiter.IsBlocked(ctx, key)
			if err != nil {
				t.Errorf("IsBlocked() error = %v", err)
				return
			}
			if blocked != tt.wantBlocked {
				t.Errorf("IsBlocked() = %v, want %v", blocked, tt.wantBlocked)
			}
		})
	}
}
