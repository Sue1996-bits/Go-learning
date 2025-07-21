package db

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm/logger"

	"gorm.io/gorm"
)

//2.0数据库初始化：连接、配置、迁移、索引

// 重写config（GORM）
type Config struct {
	DatabasePath    string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	LogLevel        logger.LogLevel //日志等级
}

func InitDatabase(config *Config) (*gorm.DB, error) {
	// gorm.Open 函数来打开一个 SQLite 数据库文件,并传入一个 gorm.Config 结构体来配置一些选项。
	db, err := gorm.Open(sqlite.Open(config.DatabasePath), &gorm.Config{
		Logger: logger.Default.LogMode(config.LogLevel),
	}) //!!!注意这里是:= /  = 即赋值而非新建！！！
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 连接池配置
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	// 执行迁移
	if err := MigrateWatermarkLog(db); err != nil {
		return nil, err
	}

	if err := CreateIndexes(db); err != nil {
		return nil, err
	}

	return db, nil

}
