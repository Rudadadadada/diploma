package redis

import (
	"time"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

func New() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",              
        DB:       0,               
    })

    _, err := RedisClient.Ping().Result()
    if err != nil {
        panic(err)
    }
}

func SetKeyWithTTL(key string, value string, ttl time.Duration) error {
    err := RedisClient.Set(key, value, ttl).Err()
    if err != nil {
        return err
    }
    return nil
}