package data

import (
	"time"
	"web-widgets/kanban-go/common"

	"gorm.io/gorm"
)

type CardUpdate struct {
	CardPosUpdate
	Card struct {
		Name         string          `json:"label"`
		Details      string          `json:"description"`
		Priority     common.FuzzyInt `json:"priority"`
		StartDate    *time.Time      `json:"start_date"`
		EndDate      *time.Time      `json:"end_date"`
		Progress     common.FuzzyInt `json:"progress"`
		Color        string          `json:"color"`
		OwnerID      common.FuzzyInt `json:"owner"`
		AttachedData []*BinaryData   `json:"attached"`
		Users        []int           `json:"users"`
	} `json:"card"`
}

type CardPosUpdate struct {
	Before   common.FuzzyInt `json:"before"`
	ColumnID common.FuzzyInt `json:"columnId"`
	RowID    common.FuzzyInt `json:"rowId"`
}

func NewCardsDAO(db *gorm.DB) *CardsDAO {
	return &CardsDAO{db}
}

type CardsDAO struct {
	db *gorm.DB
}

func (m *CardsDAO) GetAll() ([]Card, error) {
	cards := make([]Card, 0)
	err := m.db.
		Preload("AttachedData", func(db *gorm.DB) *gorm.DB {
			return m.db.Order("binary_data.id ASC")
		}).
		Preload("AssignedUsers").
		Order("`index` asc").
		Find(&cards).Error

	for i, c := range cards {
		cards[i].AssignedUsersIDs = getIDs(c.AssignedUsers)
	}
	return cards, err
}

func (m *CardsDAO) GetOne(id int) (*Card, error) {
	card := Card{}
	err := m.db.
		Preload("AttachedData", func(db *gorm.DB) *gorm.DB {
			return m.db.Order("binary_data.id ASC")
		}).
		Preload("AssignedUsers").
		First(&card, id).Error

	return &card, err
}

func (m *CardsDAO) Delete(id int) error {
	err := m.db.Exec("DELETE FROM assigned_users WHERE card_id = ?", id).Error
	if err == nil {
		err = m.db.Delete(&Card{}, id).Error
	}
	return err
}

func (m *CardsDAO) Update(id int, upd CardUpdate) error {
	c := Card{}
	err := m.db.Find(&c, id).Error
	if err != nil || c.ID == 0 {
		return err
	}

	info := upd.Card

	c.Name = info.Name
	c.Details = info.Details
	c.Priority = int(info.Priority)
	c.StartDate = info.StartDate
	c.EndDate = info.EndDate
	c.Progress = int(info.Progress)
	c.Color = info.Color
	c.AttachedData = nil
	c.AssignedUsers = nil

	err = m.db.Model(&c).Association("AssignedUsers").Clear()
	if err != nil {
		return err
	}
	if len(info.Users) > 0 {
		users := make([]User, 0)
		err := m.db.Where("id IN(?)", info.Users).Find(&users).Error
		if err != nil {
			return err
		}
		c.AssignedUsers = users
	}

	err = m.db.Save(&c).Error
	if err == nil {
		// [DIRTY] need to ensure that info.AttachedData has valid IDs
		err = m.db.Model(&BinaryData{}).Where("card_id = ?", c.ID).Update("card_id", nil).Error
		if err != nil {
			return err
		}

		if len(info.AttachedData) > 0 {
			tempIDs := make([]int, len(info.AttachedData))
			var coverId int
			for i, x := range info.AttachedData {
				tempIDs[i] = x.ID
				if x.IsCover {
					coverId = x.ID
				}
			}
			err = m.db.Model(&BinaryData{}).Where("id in (?)", tempIDs).Updates(map[string]interface{}{"card_id": c.ID, "is_cover": 0}).Error
			if err != nil {
				return err
			}
			if coverId != 0 {
				err = m.db.Model(&BinaryData{}).Where("id = ?", coverId).Update("is_cover", true).Error
			}
		}
	}

	return err
}

func (m *CardsDAO) Add(info CardUpdate) (int, error) {
	column := int(info.ColumnID)
	row := int(info.RowID)

	// get index after last item o`n the stage
	toIndex, err := m.getMaxIndex(column, row)
	if err != nil {
		return 0, err
	}

	c := Card{
		ColumnID: column,
		RowID:    row,
		Index:    toIndex,
		Name:     info.Card.Name,
	}

	err = m.db.Save(&c).Error
	return c.ID, err
}

func (m *CardsDAO) getMaxIndex(column int, row int) (int, error) {
	if column == 0 && row == 0 {
		return 0, nil
	}

	c2 := Card{}
	stm := m.db
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

func (m *CardsDAO) Move(id int, upd CardPosUpdate) error {
	c := Card{}
	err := m.db.Find(&c, id).Error
	if err != nil || c.ID == 0 {
		return err
	}

	var toIndex int
	column := int(upd.ColumnID)
	row := int(upd.RowID)
	fromIndex := c.Index
	fromColumn := c.ColumnID
	fromRow := c.RowID

	if upd.Before != 0 {
		// get move-before item
		c2 := Card{}
		err = m.db.Find(&c2, upd.Before).Error
		toIndex = c2.Index
	} else {
		// get index after last item on the stage
		toIndex, err = m.getMaxIndex(column, row)
	}
	if err != nil {
		return err
	}

	// remove item from original stage
	err = m.db.Exec("update cards set `index` = `index` - 1 where column_id = ? and row_id = ? and `index` > ?", fromColumn, fromRow, fromIndex).Error
	if err != nil {
		return err
	}

	// create place in target stage
	err = m.db.Exec("update cards set `index` = `index` + 1 where column_id = ? and row_id = ? and `index` >= ?", column, row, toIndex).Error
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
	err = m.db.Save(&c).Error

	return err
}

func getIDs(users []User) []int {
	ids := make([]int, len(users))
	for i, card := range users {
		ids[i] = card.ID
	}
	return ids
}
