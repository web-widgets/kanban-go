package main

import "github.com/jinzhu/gorm"

type CardsManager struct{}

func (m *CardsManager) GetAll() ([]Card, error) {
	cards := make([]Card, 0)
	err := db.Preload("AttachedData", func(db *gorm.DB) *gorm.DB {
		return db.Order("binary_data.id ASC")
	}).Find(&cards).Error
	return cards, err
}

func (m *CardsManager) Delete(id int) error {
	err := db.Delete(&Card{}, id).Error
	return err
}

func (m *CardsManager) Update(id int, info CardUpdate) error {
	c := Card{}
	err := db.Find(&c, id).Error
	if err != nil || c.ID == 0 {
		return err
	}

	c.Name = info.Name
	c.StageID = info.StageID

	err = db.Save(&c).Error
	if err == nil {
		// [DIRTY] need to ensure that info.AttachedData has valid IDs
		// [CRITICAL] need to ensure that only IDs are updated
		err = db.Model(&c).Association("AttachedData").Replace(info.AttachedData).Error
	}

	return err
}

func (m *CardsManager) Add(info CardUpdate) (int, error) {
	c := Card{
		Name: info.Name,
	}

	err := db.Save(&c).Error
	return c.ID, err
}
