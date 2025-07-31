package models

import (
	"survielx-backend/utility"
	"time"

	"gorm.io/gorm"
)

type Guest struct {
	ID           string         `gorm:"type:uuid;primary_key;"`
	PlateNumber  string         `json:"plate_number" gorm:"unique"`
	Model        string         `json:"model"`
	Color        string         `json:"color"`
	Type         string         `json:"type" validate:"oneof=bus car bike"`
	RegisteredBy string         `json:"registered_by" gorm:"type:uuid;"`
	User         User           `json:"user" gorm:"foreignKey:RegisteredBy"`
	CreatedAt    time.Time      `json:"createdAt"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deletedAt"`
}

func (guest *Guest) BeforeCreate(tx *gorm.DB) (err error) {
	guest.ID = utility.GenerateUUID()
	return
}
