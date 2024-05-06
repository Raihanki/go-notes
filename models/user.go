package models

import "time"

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"varchar(150)" json:"name"`
	Email     string    `gorm:"varchar(150)" json:"email"`
	Password  string    `gorm:"varchar(150)" json:"password"`
	Note      []Note    `gorm:"foreignkey:UserID" json:"note"`
	CreatedAt time.Time `json:"created_at" gorm:"timestamp"`
	UpdatedAt time.Time `json:"updated_at" gorm:"timestamp"`
}
