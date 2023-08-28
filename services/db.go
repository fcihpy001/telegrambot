package services

import (
	"context"
	"telegramBot/model"
	"telegramBot/utils"

	"gorm.io/driver/postgres"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

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
	db, err = gorm.Open(postgres.Open(utils.Config.DatabaseURL), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	logger.Info().Msg("数据库初始化成功...")
	createTable()
}

func createTable() {
	if err := db.AutoMigrate(&model.StatCount{}); err != nil {
		logger.Error().Stack().Err(err)
	}
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
	logger.Info().Msg("Redis 连接成功")
}
