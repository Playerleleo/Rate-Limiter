package main

import (
	"fmt"
	"log"
	"net/http"

	"rate-limiter/internal/config"
	"rate-limiter/internal/limiter"
	"rate-limiter/internal/middleware"
	"rate-limiter/internal/storage/redis"
)

func main() {
	// Carrega as configurações
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Inicializa o storage Redis
	storage, err := redis.NewRedisStorage(cfg.RedisHost, cfg.RedisPort, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
	}
	defer storage.Close()

	// Cria os rate limiters
	ipLimiter := limiter.NewRateLimiter(storage, limiter.Config{
		RequestsPerSecond: cfg.IPRequestsPerSecond,
		BlockDuration:     cfg.IPBlockDuration,
	})

	tokenLimiter := limiter.NewRateLimiter(storage, limiter.Config{
		RequestsPerSecond: cfg.TokenRequestsPerSecond,
		BlockDuration:     cfg.TokenBlockDuration,
	})

	// Cria o middleware
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(ipLimiter, tokenLimiter)

	// Cria o handler de exemplo
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Requisição aceita!")
	})

	// Aplica o middleware
	http.Handle("/", rateLimitMiddleware.Middleware(handler))

	// Inicia o servidor
	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Servidor iniciado na porta %s", cfg.ServerPort)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
