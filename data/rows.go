package data

import "gorm.io/gorm"

type RowUpdate struct {
	UPDBase
	Row struct {
		Name      string      `json:"label"`
		Collapsed bool        `json:"collapsed"`
	} `json:"row"`
}

type RowMove struct {
	UPDBase
	Before int `json:"before"`
}

func NewRowsDAO(db *gorm.DB) *RowsDAO {
	return &RowsDAO{db}
}

type RowsDAO struct {
	db *gorm.DB
}

func (m *RowsDAO) GetAll() ([]Row, error) {
	rows := make([]Row, 0)
	err := m.db.Order("`index` asc").Find(&rows).Error
	return rows, err
}

func (m *RowsDAO) GetOne(id int) (*Row, error) {
	r := Row{}
	err := m.db.Find(&r, id).Error
	return &r, err
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

	if info.Row.Name == "" {
		c.Collapsed = info.Row.Collapsed
	} else {
		c.Name = info.Row.Name
	}

	return m.db.Save(&c).Error
}

func (m *RowsDAO) Add(info RowUpdate) (int, error) {
	if info.RestoreID != 0 {
		err := m.db.Unscoped().Model(&Row{}).Where("id = ?", info.RestoreID).Update("deleted_at", nil).Error
		return int(info.RestoreID), err
	}

	// get index after last item o`n the stage
	toIndex, err := m.getMaxIndex()
	if err != nil {
		return 0, err
	}

	c := Row{
		Name:  info.Row.Name,
		Index: toIndex,
	}

	err = m.db.Save(&c).Error
	return c.ID, err
}

func (m *RowsDAO) getMaxIndex() (int, error) {
	r := Row{}

	err := m.db.Order("`index` desc").Take(&r).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return r.Index + 1, err
}

func (m *RowsDAO) Move(id int, before int) error {
	row := Row{}
	err := m.db.Find(&row, id).Error
	if err != nil || row.ID == 0 {
		return err
	}

	fromIndex := row.Index
	var toIndex int

	if before != 0 {
		rowBefore := Row{}
		err = m.db.Find(&rowBefore, before).Error
		toIndex = rowBefore.Index
	} else {
		// get index after last item on the stage
		toIndex, err = m.getMaxIndex()
	}
	if err != nil {
		return err
	}

	// remove item from original stage
	err = m.db.Exec("update rows set `index` = `index` - 1 where `index` > ?", fromIndex).Error
	if err != nil {
		return err
	}
	// correct index when moving from top to bottom
	if fromIndex < toIndex {
		toIndex -= 1
	}
	// create place in target stage
	err = m.db.Exec("update rows set `index` = `index` + 1 where `index` >= ?", toIndex).Error
	if err != nil {
		return err
	}

	// set item in place
	row.Index = toIndex

	err = m.db.Save(&row).Error

	return err
}
