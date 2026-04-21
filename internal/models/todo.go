package models

import "time"

type Todo struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Status      string    `json:"status" gorm:"type:enum('pending','in_progress','done');default:'pending'"`
	UserID      int       `json:"user_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Todo) TableName() string {
	return "todos"
}

type TodoQuery struct {
	Status string
	Search string
	Sort   string
	Order  string
}
