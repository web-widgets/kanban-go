package data

import (
	"time"

	"gorm.io/gorm"
)

type Board struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Card struct {
	ID      int    `json:"id"`
	Name    string `json:"label"`
	Details string `json:"description"`

	ColumnID int `json:"column"`
	RowID    int `json:"row"`

	StartDate    *time.Time    `json:"start_date"`
	EndDate      *time.Time    `json:"end_date"`
	Progress     int           `json:"progress"`
	AttachedData []*BinaryData `json:"attached"`
	Color        string        `json:"color"`
	Priority     int           `json:"priority"`

	Index            int    `json:"-"`
	AssignedUsers    []User `gorm:"many2many:assigned_users;" json:"-"`
	AssignedUsersIDs []int  `gorm:"-" json:"users"`

	DeletedAt gorm.DeletedAt `json:"-"`

	BoardID int    `json:"-"`
	Board   *Board `json:"-"`
}

type User struct {
	ID        int            `json:"id"`
	Name      string         `json:"label"`
	Avatar    string         `json:"avatar"`
	DeletedAt gorm.DeletedAt `json:"-"`

	AssignedCards []Card `gorm:"many2many:assigned_users;" json:"-"`
}

type AssignedUser struct {
	UserID    int `gorm:"primaryKey"`
	CardID    int `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt
}

type Status struct {
	ID   int    `json:"id"`
	Name string `json:"label"`

	BoardID int    `json:"-"`
	Board   *Board `json:"-"`
}

type Column struct {
	ID        int    `json:"id"`
	Name      string `json:"label"`
	Collapsed bool   `json:"collapsed"`

	Index int `json:"-"`

	DeletedAt gorm.DeletedAt `json:"-"`

	BoardID int    `json:"-"`
	Board   *Board `json:"-"`
}

type Row struct {
	ID        int    `json:"id"`
	Name      string `json:"label"`
	Collapsed bool   `json:"collapsed"`

	Index int `json:"-"`

	DeletedAt gorm.DeletedAt `json:"-"`

	BoardID int    `json:"-"`
	Board   *Board `json:"-"`
}

type BinaryData struct {
	ID      int    `json:"id"`
	Path    string `json:"-"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	IsCover bool   `json:"isCover"`

	CardID  int    `json:"-"`
	Card    *Card  `json:"-"`
	BoardID int    `json:"-"`
	Board   *Board `json:"-"`
}
