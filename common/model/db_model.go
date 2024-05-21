package model

import (
	"gorm.io/gorm"
	"time"
)

type BaseData struct {
	// ID
	ID        int64          `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"softDelete:flag" json:"deleted_at"`
}
