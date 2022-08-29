package data

import (
	"web-widgets/kanban-go/common"

	"gorm.io/gorm"
)

type ColumnUpdate struct {
	Column struct {
		ID        interface{} `json:"id"`
		Name      string      `json:"label"`
		Collapsed bool        `json:"collapsed"`
	} `json:"column"`
}

type ColumnMove struct {
	Before common.FuzzyInt `json:"before"`
}

func NewColumnsDAO(db *gorm.DB) *ColumnsDAO {
	return &ColumnsDAO{db}
}

type ColumnsDAO struct {
	db *gorm.DB
}

func (m *ColumnsDAO) GetAll() ([]Column, error) {
	columns := make([]Column, 0)
	err := m.db.Order("`index` asc").Find(&columns).Error
	return columns, err
}

func (m *ColumnsDAO) GetOne(id int) (*Column, error) {
	c := Column{}
	err := m.db.Find(&c, id).Error
	return &c, err
}

func (m *ColumnsDAO) Delete(id int) error {
	err := m.db.Delete(&Column{}, id).Error
	return err
}

func (m *ColumnsDAO) Update(id int, info ColumnUpdate) error {
	c := Column{}
	err := m.db.Find(&c, id).Error
	if err != nil || c.ID == 0 {
		return err
	}

	if info.Column.Name == "" {
		c.Collapsed = info.Column.Collapsed
	} else {
		c.Name = info.Column.Name
	}

	return m.db.Save(&c).Error
}

func (m *ColumnsDAO) Add(info ColumnUpdate) (int, error) {
	if id, ok := info.Column.ID.(float64); ok {
		err := m.db.Unscoped().Model(&Column{}).Where("id = ?", id).Update("deleted_at", nil).Error
		return int(id), err
	}

	// get index after last item o`n the stage
	toIndex, err := m.getMaxIndex()
	if err != nil {
		return 0, err
	}

	c := Column{
		Name:  info.Column.Name,
		Index: toIndex,
	}

	err = m.db.Save(&c).Error
	return c.ID, err
}

func (m *ColumnsDAO) getMaxIndex() (int, error) {
	c := Column{}

	err := m.db.Order("`index` desc").Take(&c).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return c.Index + 1, err
}

func (m *ColumnsDAO) Move(id int, before int) error {
	col := Column{}
	err := m.db.Find(&col, id).Error
	if err != nil || col.ID == 0 {
		return err
	}

	fromIndex := col.Index
	var toIndex int

	if before != 0 {
		colBefore := Column{}
		err = m.db.Find(&colBefore, before).Error
		toIndex = colBefore.Index
	} else {
		// get index after last item on the stage
		toIndex, err = m.getMaxIndex()
	}
	if err != nil {
		return err
	}

	// remove item from original stage
	err = m.db.Exec("update columns set `index` = `index` - 1 where `index` > ?", fromIndex).Error
	if err != nil {
		return err
	}
	// from right to left
	if fromIndex < toIndex {
		toIndex -= 1
	}
	// create place in target stage
	err = m.db.Exec("update columns set `index` = `index` + 1 where `index` >= ?", toIndex).Error
	if err != nil {
		return err
	}

	// set item in place
	col.Index = toIndex

	err = m.db.Save(&col).Error

	return err
}
