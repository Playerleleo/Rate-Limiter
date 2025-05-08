package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisHost              string
	RedisPort              string
	RedisPassword          string
	RedisDB                int
	IPRequestsPerSecond    int
	IPBlockDuration        int
	TokenRequestsPerSecond int
	TokenBlockDuration     int
	ServerPort             string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		return nil, err
	}

	ipRequests, err := strconv.Atoi(getEnv("IP_REQUESTS_PER_SECOND", "5"))
	if err != nil {
		return nil, err
	}

	ipBlockDuration, err := strconv.Atoi(getEnv("IP_BLOCK_DURATION_SECONDS", "300"))
	if err != nil {
		return nil, err
	}

	tokenRequests, err := strconv.Atoi(getEnv("TOKEN_REQUESTS_PER_SECOND", "10"))
	if err != nil {
		return nil, err
	}

	tokenBlockDuration, err := strconv.Atoi(getEnv("TOKEN_BLOCK_DURATION_SECONDS", "300"))
	if err != nil {
		return nil, err
	}

	return &Config{
		RedisHost:              getEnv("REDIS_HOST", "localhost"),
		RedisPort:              getEnv("REDIS_PORT", "6379"),
		RedisPassword:          getEnv("REDIS_PASSWORD", ""),
		RedisDB:                redisDB,
		IPRequestsPerSecond:    ipRequests,
		IPBlockDuration:        ipBlockDuration,
		TokenRequestsPerSecond: tokenRequests,
		TokenBlockDuration:     tokenBlockDuration,
		ServerPort:             getEnv("SERVER_PORT", "8080"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
