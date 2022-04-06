package storage

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// type Cacher interface {
// 	Get(string) Sites
// 	Set(site Sites, ttl time.Duration)
// }

type RedisStorage struct {
	Client *redis.Client
	TTL    time.Duration
	Logger *zap.Logger
}

type RedisConfig struct {
	Host string
	Port string
}

func NewRedisStorage(cfg RedisConfig, ttl time.Duration, logger *zap.Logger) (*RedisStorage, error) {
	redisLogger := logger.With(
		zap.String("storage", "Redis"),
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port))

	cli, err := NewRedisClient(cfg, redisLogger)
	if err != nil {
		redisLogger.Error("new Redis client error:", zap.Error(err))
		return nil, fmt.Errorf("new Redis client error: %s", err)
	}
	return &RedisStorage{
		Client: cli,
		TTL:    ttl,
		Logger: redisLogger,
	}, nil
}

func NewRedisClient(cfg RedisConfig, logger *zap.Logger) (*redis.Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: "",
		DB:       0,
	})

	pong, err := redisClient.Ping().Result()
	if err != nil {
		logger.Error("redis connection error", zap.Error(err))
		return nil, err
	}

	logger.Info("Redis connect", zap.String("pong", pong))

	return redisClient, nil
}

func (c *RedisStorage) Get(siteName string) (time.Duration, bool) {
	strDuration, err := c.Client.Get(siteName).Result()
	if err == redis.Nil {
		return 0, false
	}

	timeDuration, err := strconv.Atoi(strDuration)
	if err != nil {
		c.Logger.Error("Convert value to integer error", zap.Error(err))
		return 0, false
	}
	return time.Duration(timeDuration), true
}

func (c *RedisStorage) Set(searchUrlname string, duration time.Duration) {
	err := c.Client.Set(searchUrlname, duration, c.TTL).Err()
	if err != nil {
		c.Logger.Error("redis set new value error", zap.Error(err))
	}
}
