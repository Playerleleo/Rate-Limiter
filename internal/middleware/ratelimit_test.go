package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"rate-limiter/internal/limiter"
)

func TestRateLimitMiddleware(t *testing.T) {
	// Cria um mock storage
	storage := &limiter.MockStorage{
		Counts: make(map[string]int64),
	}

	// Cria os limiters
	ipLimiter := limiter.NewRateLimiter(storage, limiter.Config{
		RequestsPerSecond: 5,
		BlockDuration:     1,
	})

	tokenLimiter := limiter.NewRateLimiter(storage, limiter.Config{
		RequestsPerSecond: 10,
		BlockDuration:     1,
	})

	// Cria o middleware
	middleware := NewRateLimitMiddleware(ipLimiter, tokenLimiter)

	// Handler de teste
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	tests := []struct {
		name           string
		header         map[string]string
		requests       int
		expectedStatus int
	}{
		{
			name:           "Limite por IP",
			header:         map[string]string{},
			requests:       6,
			expectedStatus: http.StatusTooManyRequests,
		},
		{
			name:           "Limite por Token",
			header:         map[string]string{"API_KEY": "test-token"},
			requests:       11,
			expectedStatus: http.StatusTooManyRequests,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.requests; i++ {
				req := httptest.NewRequest("GET", "/", nil)
				for k, v := range tt.header {
					req.Header.Set(k, v)
				}

				rec := httptest.NewRecorder()
				middleware.Middleware(handler).ServeHTTP(rec, req)

				// Se for a última requisição, verifica o status
				if i == tt.requests-1 {
					if rec.Code != tt.expectedStatus {
						t.Errorf("Status = %v, want %v", rec.Code, tt.expectedStatus)
					}
				} else if rec.Code != http.StatusOK {
					t.Errorf("Status = %v, want %v", rec.Code, http.StatusOK)
				}
			}
		})
	}
}
