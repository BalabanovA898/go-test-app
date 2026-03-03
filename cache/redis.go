package cache

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var (
	client *redis.Client
	ttl    = 60 * time.Second
)

func init() {
	addr := os.Getenv("redis_addr")
	if addr == "" {
		log.Info("[REDIS] redis_addr not set; cache disabled")
		return
	}

	password := os.Getenv("redis_pass")
	db := parseEnvInt("redis_db", 0)

	if secStr := os.Getenv("redis_ttl_seconds"); secStr != "" {
		if sec, err := strconv.Atoi(secStr); err != nil {
			log.Warnf("[REDIS] invalid redis_ttl_seconds=%q; using %s", secStr, ttl.String())
		} else if sec <= 0 {
			log.Warnf("[REDIS] redis_ttl_seconds=%d; using %s", sec, ttl.String())
		} else {
			ttl = time.Duration(sec) * time.Second
		}
	}

	c := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := c.Ping(ctx).Err(); err != nil {
		log.Warnf("[REDIS] unable to connect to %s: %v; cache disabled", addr, err)
		_ = c.Close()
		return
	}

	client = c
	log.Info("[REDIS] Connected on " + addr)
}

func GetRedis() *redis.Client {
	return client
}

func TTL() time.Duration {
	return ttl
}

func NoteKey(id string) string {
	return "note:" + id
}

func parseEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Warnf("[%s] invalid int value %q; using %d", key, v, fallback)
		return fallback
	}
	return n
}
