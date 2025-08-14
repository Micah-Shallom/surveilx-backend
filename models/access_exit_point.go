package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type AccessExitPoint struct {
    ID        string         `gorm:"column:id;type:uuid;primaryKey;" json:"id"`
    Name      string         `json:"name" gorm:"column:name;unique"`
    CreatedAt time.Time      `json:"createdAt" gorm:"column:created_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt" gorm:"column:deleted_at"`
}

func (point *AccessExitPoint) BeforeCreate(tx *gorm.DB) (err error) {
	point.ID = utility.GenerateUUID()
	return
}
