package main

import "gorm.io/gorm"

type CardsManager struct{}

func (m *CardsManager) GetAll() ([]Card, error) {
	cards := make([]Card, 0)
	err := db.Preload("AttachedData", func(db *gorm.DB) *gorm.DB {
		return db.Order("binary_data.id ASC")
	}).Order("`index` asc").Find(&cards).Error
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
	c.StageID = int(info.StageID)
	c.OwnerID = int(info.OwnerID)
	c.Details = info.Details
	c.StartDate = info.StartDate
	c.AttachedData = nil

	err = db.Save(&c).Error
	if err == nil {
		// [DIRTY] need to ensure that info.AttachedData has valid IDs
		// [CRITICAL] need to ensure that only IDs are updated
		err = db.Model(&c).Association("AttachedData").Replace(info.AttachedData)
	}

	return err
}

func (m *CardsManager) Add(info CardUpdate) (int, error) {
	toStage := int(info.StageID)

	// get index after last item on the stage
	toIndex, err := m.getMaxIndex(toStage)
	if err != nil {
		return 0, err
	}

	c := Card{
		StageID: toStage,
		Index:   toIndex,
		Name:    info.Name,
	}

	err = db.Save(&c).Error
	return c.ID, err
}

func (m *CardsManager) getMaxIndex(stage int) (int, error) {
	c2 := Card{}
	err := db.Where("stage_id=?", stage).Order("`index` desc").Take(&c2).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return c2.Index + 1, err
}

func (m *CardsManager) Move(id int, card CardUpdate, before FuzzyInt) error {
	c := Card{}
	err := db.Find(&c, id).Error
	if err != nil || c.ID == 0 {
		return err
	}

	var toIndex int
	toStage := int(card.StageID)
	fromIndex := c.Index
	fromStage := c.StageID

	if before != 0 {
		// get move-before item
		c2 := Card{}
		err = db.Find(&c2, before).Error
		toIndex = c2.Index
	} else {
		// get index after last item on the stage
		toIndex, err = m.getMaxIndex(toStage)
	}
	if err != nil {
		return err
	}

	// remove item from original stage
	err = db.Exec("update cards set `index` = `index` - 1 where stage_id = ? and `index` > ?", fromStage, fromIndex).Error
	if err != nil {
		return err
	}

	// create place in target stage
	err = db.Exec("update cards set `index` = `index` + 1 where stage_id = ? and `index` >= ?", toStage, toIndex).Error
	if err != nil {
		return err
	}

	// set item in place
	c.Index = toIndex
	// correct index when moving from top to bottom in the same list
	if fromStage == toStage && fromIndex < toIndex {
		c.Index -= 1
	}
	c.StageID = toStage
	err = db.Save(&c).Error

	return err
}
