package database

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func RedisClient(dbno int) *redis.Client {
	rd := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDRESS"),
		Password: os.Getenv("DB_PASSWORD"),
		DB:       dbno,
	})
	return rd
}
