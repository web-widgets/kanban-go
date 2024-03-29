package data

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Debug = 1
var features *Features

func logError(e error) {
	if e != nil && Debug > 0 {
		log.Printf("[ERROR]\n%s\n", e)
	}
}

type DBConfig struct {
	Path         string
	ResetOnStart bool
}

type Features struct {
	WithVotes    bool `yaml:"withVotes"`
	WithComments bool `yaml:"withComments"`
}

type DAO struct {
	db *gorm.DB

	Cards    *CardsDAO
	Rows     *RowsDAO
	Columns  *ColumnsDAO
	Files    *FilesDAO
	Users    *UsersDAO
	Comments *CommentsDAO
	Links    *LinksDAO
}

func (d *DAO) GetDB() *gorm.DB {
	return d.db
}

func (d *DAO) mustExec(sql string) {
	err := d.db.Exec(sql).Error
	if err != nil {
		panic(err)
	}
}

func NewDAO(config DBConfig, url, drive string, featureCfg Features) *DAO {
	db, err := gorm.Open(sqlite.Open(config.Path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic("failed to connect database")
	}

	features = &featureCfg

	db.AutoMigrate(&Card{})
	db.AutoMigrate(&Column{})
	db.AutoMigrate(&Row{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&AssignedUser{})
	db.AutoMigrate(&Status{})
	db.AutoMigrate(&BinaryData{})
	db.AutoMigrate(&Vote{})
	db.AutoMigrate(&Comment{})
	db.AutoMigrate(&Link{})

	dao := &DAO{
		db:       db,
		Cards:    NewCardsDAO(db),
		Rows:     NewRowsDAO(db),
		Columns:  NewColumnsDAO(db),
		Files:    NewFilesDAO(db, url, drive),
		Users:    NewUsersDAO(db),
		Comments: NewCommentsDAO(db),
		Links:    NewLinksDAO(db),
	}

	if config.ResetOnStart {
		dataDown(dao)
		dataUp(dao)
	}

	return dao
}
