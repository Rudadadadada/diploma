package redis

import (
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