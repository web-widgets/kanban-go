package data

type RowUpdate struct {
	Name string `json:"label"`
}

type RowsManager struct{}

func (m *RowsManager) GetAll() ([]Row, error) {
	rows := make([]Row, 0)
	err := db.Find(&rows).Error
	return rows, err
}

func (m *RowsManager) Delete(id int) error {
	err := db.Delete(&Row{}, id).Error
	return err
}

func (m *RowsManager) Update(id int, info RowUpdate) error {
	c := Row{}
	err := db.Find(&c, id).Error
	if err != nil || c.ID == 0 {
		return err
	}

	c.Name = info.Name

	return db.Save(&c).Error
}

func (m *RowsManager) Add(info RowUpdate) (int, error) {
	c := Row{
		Name: info.Name,
	}

	err := db.Save(&c).Error
	return c.ID, err
}
