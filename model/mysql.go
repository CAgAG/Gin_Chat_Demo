package model

import (
	logging "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func NewDatabase(connString string) {
	// dsn := "test_root:123456@tcp(127.0.0.1:3306)/gin_test?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := connString
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "chat_demo_",
			SingularTable: true,
		},
		// 打印对应的 sql 语句
		// Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	db.AutoMigrate(&User{}) // 自动同步
	DB = db
}
