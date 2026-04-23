package models

type Role struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement;not null"`
	RoleName string `json:"role_name" gorm:"unique;type:varchar(50);not null"`
}

func (Role) TableName() string {
	return "roles"
}
