package data

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CommentUpdate struct {
	Text     string     `json:"text"`
	PostedAt *time.Time `json:"date"`
}

type CommentsDAO struct {
	db *gorm.DB
}

func NewCommentsDAO(db *gorm.DB) *CommentsDAO {
	return &CommentsDAO{db}
}

func (d *CommentsDAO) GetOne(id int) (Comment, error) {
	comment := Comment{}
	err := d.db.Find(&comment, id).Error
	return comment, err
}

func (d *CommentsDAO) Add(cardId, userId int, upd CommentUpdate) (int, error) {
	comment := Comment{
		UserID:   userId,
		CardID:   cardId,
		Text:     upd.Text,
		PostedAt: upd.PostedAt,
	}
	err := d.db.Create(&comment).Error

	return comment.ID, err
}

func (d *CommentsDAO) Update(commentId, cardId, userId int, upd CommentUpdate) error {
	comment, err := d.GetOne(commentId)
	if err != nil {
		return err
	}
	if comment.UserID != userId || comment.CardID != cardId {
		return fmt.Errorf("access denied")
	}

	comment.Text = upd.Text

	err = d.db.Save(&comment).Error

	return err
}

func (d *CommentsDAO) Delete(commentId, cardId, userId int) error {
	comment, err := d.GetOne(commentId)
	if err != nil {
		return err
	}

	if comment.UserID != userId || comment.CardID != cardId {
		return fmt.Errorf("access denied")
	}

	err = d.db.Delete(&Comment{}, commentId).Error

	return err
}
