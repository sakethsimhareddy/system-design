package db

import (
	"github.com/go-redis/redis"

)

type redisDB struct{
	client *redis.Client
}

type RedisDB interface {
	Save(shortCode string,originalURL string) error
	Get(shortCode string) (*string, error)
}

func NewRedisDB() RedisDB {
	return &redisDB{
		client: redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
		}),
	}
}


func (r *redisDB) Save(shortCode string,originalURL string) error {
	return r.client.Set(shortCode, originalURL, 0).Err()
}

func (r *redisDB) Get(shortCode string) (*string, error) {
	originalURL, err := r.client.Get(shortCode).Result()
	if err != nil {
		return nil, err	
	}
	return &originalURL, nil
}
