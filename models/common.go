package models

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	Id        uint           `gorm:"primaryKey;type:int(11);unique" json:"id"`
	DeletedAt gorm.DeletedAt `gorm:"type:datetime" json:"deleted_at"`
	CreatedAt time.Time      `gorm:"type:datetime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:datetime" json:"updated_at"`
}
