package main

type StagesManager struct{}

func (m *StagesManager) GetAll() ([]Stage, error) {
	cards := make([]Stage, 0)
	err := db.Find(&cards).Error
	return cards, err
}

func (m *StagesManager) Delete(id int) error {
	err := db.Delete(&Stage{}, id).Error
	return err
}

func (m *StagesManager) Update(id int, info StageUpdate) error {
	c := Stage{}
	err := db.Find(&c, id).Error
	if err != nil || c.ID == 0 {
		return err
	}

	c.Name = info.Name

	return db.Save(&c).Error
}

func (m *StagesManager) Add(info StageUpdate) (int, error) {
	c := Stage{
		Name: info.Name,
	}

	err := db.Save(&c).Error
	return c.ID, err
}
