package services

import (
	"context"
	"fmt"
	"telegramBot/model"
	"telegramBot/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
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
	// initPostgres(utils.Config.DatabaseURL)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.Config.Mysql.UserName,
		utils.Config.Mysql.Passwd,
		utils.Config.Mysql.Address,
		utils.Config.Mysql.Database)
	fmt.Println(dsn)
	InitMysql(dsn)

	initRedis(utils.Config.RedisURL)
}

//	func initPostgres(dsn string) {
//		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
func InitMysql(dsn string) {
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
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

//lint:ignore U1000 ignore unused lint
func createTable() {
	if err := db.AutoMigrate(&model.User{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.UserChat{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.UserAction{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.ChatGroup{}); err != nil {
		logger.Error().Stack().Err(err)
	}
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
