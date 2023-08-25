package services

import (
	"context"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"telegramBot/model"
	"telegramBot/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
	err error
)

func InitDB() {

	initPostgres()

	initRedis()
}

func initPostgres() {
	fmt.Println(utils.Config.DatabaseURL)
	logMode := logger.Info

	db, err = gorm.Open(postgres.Open(utils.Config.DatabaseURL), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("数据库初始化成功...")
	createTable()
}

func createTable() {
	db.AutoMigrate(&model.StatCount{})
}

func initRedis() {
	opts, err := redis.ParseURL(utils.Config.RedisURL)
	if err != nil {
		panic(err)
	}

	rdb = redis.NewClient(opts)
	if err = rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	log.Println("Redis 连接成功")
}
