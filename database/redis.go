package database

import (
	"app-inventory/utils"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// InitRedis initializes and tests Redis client connection
func InitRedis(config utils.RedisConfig, logger *zap.Logger) *redis.Client {
	if !config.Enabled {
		logger.Info("Redis caching is disabled via configuration")
		return nil
	}

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Warn("Failed to connect to Redis, proceeding with cache disabled",
			zap.String("addr", addr),
			zap.Error(err),
		)
		return nil
	}

	logger.Info("Successfully connected to Redis cache", zap.String("addr", addr))
	return rdb
}
