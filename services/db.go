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

func Init(ctx context.Context) {
	InitDB()
	go StatsRoutine(ctx)
}

func InitDB() {
	initPostgres(utils.Config.DatabaseURL)

	initRedis(utils.Config.RedisURL)
}

func initPostgres(dsn string) {
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}
	logger.Info().Msg("数据库初始化成功...")
	// createTable()
}

//lint:ignore U1000 ignore unused lint
func createTable() {
	if err := db.AutoMigrate(&model.StatCount{}); err != nil {
		logger.Error().Stack().Err(err)
	}
}

func initRedis(uri string) {
	opts, err := redis.ParseURL(uri)
	if err != nil {
		panic(err)
	}

	rdb = redis.NewClient(opts)
	if err = rdb.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	logger.Info().Msg("Redis 连接成功")
}
