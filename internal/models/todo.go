package models

import "time"

type Todo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"type:varchar(255);not null"`
	Description string    `json:"description" gorm:"type:text"`
	Status      string    `json:"status" gorm:"type:enum('pending','in_progress','done');default:'pending'"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Todo) TableName() string {
	return "todos"
}
