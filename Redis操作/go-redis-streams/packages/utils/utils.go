package utils

import (
	"github.com/go-redis/redis/v7"
)

//NewRedisClient create a new instace of client redis
func NewRedisClient(redisOp *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(redisOp)

	_, err := client.Ping().Result()
	return client, err

}
