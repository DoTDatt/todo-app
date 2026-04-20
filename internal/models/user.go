package models

import "time"

type User struct {
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	Email     string    `json:"email" gorm:"unique;type:varchar(100);not null"`
	Password  string    `json:"password" gorm:"type:varchar(255);not null"`
	RoleID    int       `json:"role_id" gorm:"column:role_id;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
