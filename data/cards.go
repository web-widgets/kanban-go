package data

import (
	"fmt"
	"time"
	"web-widgets/kanban-go/common"

	"gorm.io/gorm"
)

type CardUpdate struct {
	Meta MetaInfo `json:"$meta"`
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
		Votes        []int           `json:"votes"`
		RowID        common.FuzzyInt `json:"row"`
		ColumnID     common.FuzzyInt `json:"column"`
	} `json:"card"`
}

type CardPosUpdate struct {
	Meta     MetaInfo        `json:"$meta"`
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

	op := m.db.
		Preload("AttachedData", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("binary_data.id ASC")
		}).
		Preload("AssignedUsers").
		Order("`index` asc")

	if features.WithVotes {
		op.Preload("Votes")
	}
	if features.WithComments {
		op.Preload("Comments")
	}

	err := op.Find(&cards).Error

	for i, c := range cards {
		cards[i].AssignedUsersIDs = getUserIDs(c.AssignedUsers)
		if features.WithVotes {
			cards[i].VotesUsersIDs = getVoteUserIDs(c.Votes)
		}
	}

	return cards, err
}

func (m *CardsDAO) GetColumn(id int) ([]Card, error) {
	cards := make([]Card, 0)

	op := m.db.
		Preload("AttachedData", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("binary_data.id ASC")
		}).
		Preload("AssignedUsers").
		Order("`index` asc")

	if features.WithVotes {
		op.Preload("Votes")
	}
	if features.WithComments {
		op.Preload("Comments")
	}

	err := op.Find(&cards, "column_id = ?", id).Error

	for i, c := range cards {
		cards[i].AssignedUsersIDs = getUserIDs(c.AssignedUsers)
		if features.WithVotes {
			cards[i].VotesUsersIDs = getVoteUserIDs(c.Votes)
		}
	}

	return cards, err
}

func (m *CardsDAO) GetOne(id int) (*Card, error) {
	card := Card{}
	op := m.db.
		Preload("AttachedData", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("binary_data.id ASC")
		}).
		Preload("AssignedUsers")

	if features.WithVotes {
		op.Preload("Votes")
		card.VotesUsersIDs = getVoteUserIDs(card.Votes)
	}
	if features.WithComments {
		op.Preload("Comments")
	}

	err := op.First(&card, id).Error

	card.AssignedUsersIDs = getUserIDs(card.AssignedUsers)

	return &card, err
}

func (m *CardsDAO) Delete(id int) error {
	err := m.db.Where("card_id = ?", id).Delete(&AssignedUser{}).Error
	if err != nil {
		return err
	}

	if features.WithVotes {
		err = m.db.Where("card_id = ?", id).Delete(&Vote{}).Error
	}

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
	if info.Meta.RestoreID != 0 {
		err := m.db.Unscoped().Model(&Card{}).Where("id = ?", info.Meta.RestoreID).Update("deleted_at", nil).Error
		if err == nil {
			err = m.db.Unscoped().Model(&AssignedUser{}).Where("card_id = ?", info.Meta.RestoreID).Update("deleted_at", nil).Error
		}
		return int(info.Meta.RestoreID), err
	}

	column := int(info.Card.ColumnID)
	row := int(info.Card.RowID)

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

func (m *CardsDAO) SetVote(cid, user int) error {
	if cid == 0 {
		return fmt.Errorf("card ID not defined")
	}
	if user == 0 {
		return fmt.Errorf("user ID not defined")
	}

	vote := Vote{}
	err := m.db.Where("card_id = ? AND user_id = ?", cid, user).Find(&vote).Error
	if err != nil {
		return err
	}

	if vote.CardID != 0 && vote.UserID != 0 {
		// vote already exists
		return nil
	}

	vote = Vote{
		CardID: cid,
		UserID: user,
	}

	err = m.db.Create(&vote).Error

	return err
}

func (m *CardsDAO) RemoveVote(cid, user int) error {
	if cid == 0 {
		return fmt.Errorf("card ID not defined")
	}
	if user == 0 {
		return fmt.Errorf("user ID not defined")
	}

	return m.db.Where("card_id = ? AND user_id = ?", cid, user).Delete(&Vote{}).Error
}

func getUserIDs(users []User) []int {
	ids := make([]int, len(users))
	for i, card := range users {
		ids[i] = card.ID
	}
	return ids
}

func getVoteUserIDs(votes []Vote) []int {
	ids := make([]int, len(votes))
	for i, v := range votes {
		ids[i] = v.UserID
	}
	return ids
}
