package data

import "gorm.io/gorm"

type ColumnUpdate struct {
	Name string `json:"label"`
}

func NewColumnsDAO(db *gorm.DB) *ColumnsDAO {
	return &ColumnsDAO{db}
}

type ColumnsDAO struct {
	db *gorm.DB
}

func (m *ColumnsDAO) GetAll() ([]Column, error) {
	columns := make([]Column, 0)
	err := m.db.Find(&columns).Error
	return columns, err
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

	c.Name = info.Name

	return m.db.Save(&c).Error
}

func (m *ColumnsDAO) Add(info ColumnUpdate) (int, error) {
	c := Column{
		Name: info.Name,
	}

	err := m.db.Save(&c).Error
	return c.ID, err
}
