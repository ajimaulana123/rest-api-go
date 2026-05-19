package models

import "time"

type Item struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"size:200;not null"`
	Description string    `json:"description" gorm:"type:text"`
	UserID      uint      `json:"user_id" gorm:"not null;index:idx_items_user_created,priority:1"`
	User        User      `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	CreatedAt   time.Time `json:"created_at" gorm:"index:idx_items_user_created,priority:2,sort:desc"`
	UpdatedAt   time.Time `json:"updated_at"`
}
