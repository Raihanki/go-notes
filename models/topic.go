package models

import "time"

type Topic struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"varchar(150)" json:"name"`
	Note      []Note    `gorm:"foreignkey:TopicID" json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
