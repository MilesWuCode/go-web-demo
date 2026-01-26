package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string
	Email    string
	Password string
	Posts    []Post `gorm:"many2many:user_posts"`
}
