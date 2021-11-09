package data

import (
	"time"
	"web-widgets/kanban-go/common"

	"gorm.io/gorm"
)

type CardUpdate struct {
	Name         string          `json:"label"`
	ColumnID     common.FuzzyInt `json:"column"`
	RowID        common.FuzzyInt `json:"row"`
	Details      string          `json:"description"`
	Priority     common.FuzzyInt `json:"priority"`
	StartDate    *time.Time      `json:"start_date"`
	EndDate      *time.Time      `json:"end_date"`
	Progress     common.FuzzyInt `json:"porgress"`
	Color        string          `json:"color"`
	OwnerID      common.FuzzyInt `json:"owner"`
	AttachedData []*BinaryData   `json:"attached,omitempty"`
}

type CardMove struct {
	Card   CardUpdate      `json:"card"`
	Before common.FuzzyInt `json:"before"`
}

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
	c.Details = info.Details
	c.Priority = int(info.Priority)
	c.StartDate = info.StartDate
	c.EndDate = info.EndDate
	c.Progress = int(info.Progress)
	c.Color = info.Color
	c.AttachedData = nil

	err = db.Save(&c).Error
	if err == nil {
		// [DIRTY] need to ensure that info.AttachedData has valid IDs
		err = db.Model(&BinaryData{}).Where("card_id = ?", c.ID).Update("card_id", nil).Error
		if err != nil {
			return err
		}

		if len(info.AttachedData) > 0 {
			tempIDs := make([]int, len(info.AttachedData))
			for i, x := range info.AttachedData {
				tempIDs[i] = x.ID
			}
			err = db.Model(&BinaryData{}).Where("id in (?)", tempIDs).Update("card_id", c.ID).Error
		}
	}

	return err
}

func (m *CardsManager) Add(info CardUpdate) (int, error) {
	column := int(info.ColumnID)
	row := int(info.ColumnID)

	// get index after last item o`n the stage
	toIndex, err := m.getMaxIndex(column, row)
	if err != nil {
		return 0, err
	}

	c := Card{
		ColumnID: column,
		RowID:    row,
		Index:    toIndex,
		Name:     info.Name,
	}

	err = db.Save(&c).Error
	return c.ID, err
}

func (m *CardsManager) getMaxIndex(column int, row int) (int, error) {
	if column == 0 && row == 0 {
		return 0, nil
	}

	c2 := Card{}
	stm := db
	if column != 0 {
		stm = stm.Where("column_id=?", column)
	}
	if row != 0 {
		stm = stm.Where("row_id=?", row)
	}

	err := stm.Order("`index` desc").Take(&c2).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return c2.Index + 1, err
}

func (m *CardsManager) Move(id int, card CardUpdate, before int) error {
	c := Card{}
	err := db.Find(&c, id).Error
	if err != nil || c.ID == 0 {
		return err
	}

	var toIndex int
	column := int(card.ColumnID)
	row := int(card.RowID)
	fromIndex := c.Index
	fromColumn := c.ColumnID
	fromRow := c.RowID

	if before != 0 {
		// get move-before item
		c2 := Card{}
		err = db.Find(&c2, before).Error
		toIndex = c2.Index
	} else {
		// get index after last item on the stage
		toIndex, err = m.getMaxIndex(column, row)
	}
	if err != nil {
		return err
	}

	// remove item from original stage
	err = db.Exec("update cards set `index` = `index` - 1 where column_id = ? and row_id = ? and `index` > ?", fromColumn, fromRow, fromIndex).Error
	if err != nil {
		return err
	}

	// create place in target stage
	err = db.Exec("update cards set `index` = `index` + 1 where column_id = ? and row_id = ? and `index` >= ?", column, row, toIndex).Error
	if err != nil {
		return err
	}

	// set item in place
	c.Index = toIndex
	// correct index when moving from top to bottom in the same list
	if fromColumn == column && fromRow == row && fromIndex < toIndex {
		c.Index -= 1
	}

	c.ColumnID = column
	c.RowID = row
	err = db.Save(&c).Error

	return err
}
