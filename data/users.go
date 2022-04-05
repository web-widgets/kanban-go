package data

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"

	"gorm.io/gorm"
)

func NewUsersDAO(db *gorm.DB) *UsersDAO {
	return &UsersDAO{db}
}

type UsersDAO struct {
	db *gorm.DB
}

func (m *UsersDAO) GetAll() ([]User, error) {
	users := make([]User, 0)
	err := m.db.Find(&users).Error

	for i, u := range users {
		r, err := getAvatar(&u)
		if err != nil {
			return nil, err
		}
		users[i].Avatar = r
	}

	return users, err
}

func getAvatar(user *User) (string, error) {
	if user.Avatar == "" {
		return "", nil
	}

	bytes, err := ioutil.ReadFile(user.Avatar)
	if err != nil {
		return "", err
	}
	base64Enc := toBase64(bytes)

	return base64Enc, nil
}

func toBase64(bytes []byte) string {
	// hardcoded .jpg format
	mimeType := "image/jpg"
	s := base64.StdEncoding.EncodeToString(bytes)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, s)
}
