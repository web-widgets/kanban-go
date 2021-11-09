package data

type ColumnUpdate struct {
	Name string `json:"label"`
}

type ColumnsManager struct{}

func (m *ColumnsManager) GetAll() ([]Column, error) {
	columns := make([]Column, 0)
	err := db.Find(&columns).Error
	return columns, err
}

func (m *ColumnsManager) Delete(id int) error {
	err := db.Delete(&Column{}, id).Error
	return err
}

func (m *ColumnsManager) Update(id int, info ColumnUpdate) error {
	c := Column{}
	err := db.Find(&c, id).Error
	if err != nil || c.ID == 0 {
		return err
	}

	c.Name = info.Name

	return db.Save(&c).Error
}

func (m *ColumnsManager) Add(info ColumnUpdate) (int, error) {
	c := Column{
		Name: info.Name,
	}

	err := db.Save(&c).Error
	return c.ID, err
}
