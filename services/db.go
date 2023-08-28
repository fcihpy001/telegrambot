package services

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"telegramBot/model"
	"telegramBot/utils"

	"gorm.io/gorm"
)

var (
	db  *gorm.DB
	rdb *redis.Client
	err error
)

func InitDB() {

	//InitMysql()

	initRedis()
}

func InitMysql() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.Config.Mysql.UserName,
		utils.Config.Mysql.Passwd,
		utils.Config.Mysql.Address,
		utils.Config.Mysql.Database)
	fmt.Println(dsn)
	logMode := logger.Info

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
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
	return
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
