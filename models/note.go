package models

import "time"

type Note struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	UserID    uint      `json:"user_id"`
	User      User      `gorm:"foreignkey:UserID" json:"user"`
	Title     string    `gorm:"varchar(150)" json:"title"`
	Content   string    `gorm:"text" json:"content"`
	TopicID   uint      `json:"topic_id"`
	Topic     Topic     `gorm:"foreignkey:TopicID" json:"topic"`
	CreatedAt time.Time `gorm:"timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"timestamp" json:"updated_at"`
}
