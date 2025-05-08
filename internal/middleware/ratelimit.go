package middleware

import (
	"net/http"
	"strings"

	"rate-limiter/internal/limiter"
)

type RateLimitMiddleware struct {
	ipLimiter    *limiter.RateLimiter
	tokenLimiter *limiter.RateLimiter
}

func NewRateLimitMiddleware(ipLimiter, tokenLimiter *limiter.RateLimiter) *RateLimitMiddleware {
	return &RateLimitMiddleware{
		ipLimiter:    ipLimiter,
		tokenLimiter: tokenLimiter,
	}
}

func (m *RateLimitMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verifica o token primeiro
		token := r.Header.Get("API_KEY")
		if token != "" {
			allowed, err := m.tokenLimiter.CheckLimit(r.Context(), token)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if !allowed {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		} else {
			// Se não houver token, verifica o IP
			ip := strings.Split(r.RemoteAddr, ":")[0]
			allowed, err := m.ipLimiter.CheckLimit(r.Context(), ip)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			if !allowed {
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
