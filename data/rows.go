package data

import "gorm.io/gorm"

type RowUpdate struct {
	Name string `json:"label"`
}

func NewRowsDAO(db *gorm.DB) *RowsDAO {
	return &RowsDAO{db}
}

type RowsDAO struct {
	db *gorm.DB
}

func (m *RowsDAO) GetAll() ([]Row, error) {
	rows := make([]Row, 0)
	err := m.db.Find(&rows).Error
	return rows, err
}

func (m *RowsDAO) Delete(id int) error {
	err := m.db.Delete(&Row{}, id).Error
	return err
}

func (m *RowsDAO) Update(id int, info RowUpdate) error {
	c := Row{}
	err := m.db.Find(&c, id).Error
	if err != nil || c.ID == 0 {
		return err
	}

	c.Name = info.Name

	return m.db.Save(&c).Error
}

func (m *RowsDAO) Add(info RowUpdate) (int, error) {
	c := Row{
		Name: info.Name,
	}

	err := m.db.Save(&c).Error
	return c.ID, err
}
