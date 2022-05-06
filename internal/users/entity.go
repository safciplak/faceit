package users

import (
	"time"
)

type User struct {
	ID        string    `gorm:"column:id"`
	FirstName string    `gorm:"column:first_name"`
	LastName  string    `gorm:"column:last_name"`
	Nickname  string    `gorm:"column:nickname"`
	Password  string    `gorm:"column:password"`
	Email     string    `gorm:"column:email"`
	Country   string    `gorm:"column:country"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

func (u User) TableName() string {
	return "users"
}
