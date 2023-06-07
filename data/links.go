package data

import (
	"web-widgets/kanban-go/common"

	"gorm.io/gorm"
)

type LinkUpdate struct {
	Meta MetaInfo `json:"$meta"`
	Link struct {
		MasterID common.FuzzyInt `json:"masterId"`
		SlaveID  common.FuzzyInt `json:"slaveId"`
		Relation string          `json:"relation"`
	} `json:"link"`
}

func NewLinksDAO(db *gorm.DB) *LinksDAO {
	return &LinksDAO{db}
}

type LinksDAO struct {
	db *gorm.DB
}

func (m *LinksDAO) GetAll() ([]Link, error) {
	links := make([]Link, 0)
	err := m.db.Order("`index` asc").Find(&links).Error
	return links, err
}

func (m *LinksDAO) GetOne(id int) (*Link, error) {
	l := Link{}
	err := m.db.Find(&l, id).Error
	return &l, err
}

func (m *LinksDAO) Delete(id int) error {
	err := m.db.Delete(&Link{}, id).Error
	return err
}

func (m *LinksDAO) Add(info LinkUpdate) (int, error) {
	if info.Meta.RestoreID != 0 {
		err := m.db.Unscoped().Model(&Link{}).Where("id = ?", info.Meta.RestoreID).Update("deleted_at", nil).Error
		return int(info.Meta.RestoreID), err
	}

	// get index after last item o`n the stage
	toIndex, err := m.getMaxIndex()
	if err != nil {
		return 0, err
	}

	c := Link{
		MasterID: int(info.Link.MasterID),
		SlaveID:  int(info.Link.SlaveID),
		Relation: info.Link.Relation,
		Index:    toIndex,
	}

	err = m.db.Save(&c).Error
	return c.ID, err
}

func (m *LinksDAO) getMaxIndex() (int, error) {
	l := Link{}

	err := m.db.Order("`index` desc").Take(&l).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return l.Index + 1, err
}
