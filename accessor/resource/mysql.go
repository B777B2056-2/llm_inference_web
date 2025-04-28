package resource

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"llm_online_inference/scheduler/confparser"
	"time"
)

var DB *gorm.DB

func buildDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		confparser.ResourceConfig.MySQL.Username,
		confparser.ResourceConfig.MySQL.Password,
		confparser.ResourceConfig.MySQL.Host,
		confparser.ResourceConfig.MySQL.Port,
		confparser.ResourceConfig.MySQL.DBName,
	)
}

func initMySQL() {
	db, err := gorm.Open(mysql.Open(buildDSN()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to connect database")
	}
	sqlDB.SetMaxIdleConns(confparser.ResourceConfig.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(confparser.ResourceConfig.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(
		time.Duration(
			confparser.ResourceConfig.MySQL.ConnMaxLifetimeInSecond,
		) * time.Second,
	)

	DB = db
}
