package services

import (
	"os"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
)

func Init() {
	var err error
	dburl := os.Getenv("DATABASE_URL")
	redisUrl := os.Getenv("REDIS_URL")

	db, err = OpenDB(dburl)
	if err != nil {
		panic(err)
	}
	rdb = OpenRedis(redisUrl)
}
