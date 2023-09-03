package services

import (
	"context"
	"fmt"
	"log"
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
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		utils.Config.Mysql.UserName,
		utils.Config.Mysql.Passwd,
		utils.Config.Mysql.Address,
		utils.Config.Mysql.Database)
	// fmt.Println(dsn)
	InitMysql(dsn)
	createTable()

	initRedis(utils.Config.RedisURL)
}

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
	if err := db.AutoMigrate(&model.WelcomeSetting{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.InviteSetting{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.Reply{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.ReplySetting{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.ProhibitedSetting{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.Punishment{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.NewMemberCheck{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.UserCheck{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.UserCautions{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.SpamSetting{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.FloodSetting{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.DarkModelSetting{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.Solitaire{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.SolitaireMessage{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.VerifySetting{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.ScheduleMsg{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	if err := db.AutoMigrate(&model.Schedule{}); err != nil {
		logger.Error().Stack().Err(err)
	}
	log.Println("数据表创建成功...")
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
