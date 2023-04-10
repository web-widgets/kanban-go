package data

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type CommentUpdate struct {
	CardID  int `json:"cardId"`
	Comment struct {
		Text     string     `json:"text"`
		PostedAt *time.Time `json:"date"`
	}
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

func (d *CommentsDAO) Add(userId int, upd CommentUpdate) (int, error) {
	comment := Comment{
		UserID:   userId,
		CardID:   upd.CardID,
		Text:     upd.Comment.Text,
		PostedAt: upd.Comment.PostedAt,
	}
	err := d.db.Create(&comment).Error

	return comment.ID, err
}

func (d *CommentsDAO) Update(commentId, userId int, upd CommentUpdate) error {
	comment, err := d.GetOne(commentId)
	if err != nil {
		return err
	}
	if comment.UserID != userId {
		return fmt.Errorf("access denied")
	}

	comment.Text = upd.Comment.Text

	err = d.db.Save(&comment).Error

	return err
}

func (d *CommentsDAO) Delete(commentId, userId int) error {
	comment, err := d.GetOne(commentId)
	if err != nil {
		return err
	}

	if comment.UserID != userId {
		return fmt.Errorf("access denied")
	}

	err = d.db.Delete(&Comment{}, commentId).Error

	return err
}
