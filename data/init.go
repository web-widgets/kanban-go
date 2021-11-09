package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	Path         string
	ResetOnStart bool
}

var db *gorm.DB

func Init(config DBConfig) *gorm.DB {
	var err error
	db, err = gorm.Open(sqlite.Open(config.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Card{})
	db.AutoMigrate(&Column{})
	db.AutoMigrate(&Row{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Status{})
	db.AutoMigrate(&BinaryData{})

	if config.ResetOnStart {
		dataDown()
		dataUp()
	}

	return db
}

func mustExec(sql string) {
	err := db.Exec(sql).Error
	if err != nil {
		panic(err)
	}
}
