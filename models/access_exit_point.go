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
	DeletedAt gorm.DeletedAt `json:"deletedAt" gorm:"column:deleted_at"`
}

func (point *AccessExitPoint) BeforeCreate(tx *gorm.DB) (err error) {
	point.ID = utility.GenerateUUID()
	return
}

// In models/models.go
type PendingVehicleExit struct {
	ID            string    `gorm:"column:id;type:uuid;primaryKey;" json:"id"`
	PlateNumber   string    `json:"plateNumber" gorm:"column:plate_number"`
	VehicleID     string    `json:"vehicleId" gorm:"column:vehicle_id"`
	UserID        string    `json:"userId" gorm:"column:user_id"` // Owner's user ID, for notification targeting
	ExitPointID   string    `json:"exitPointId" gorm:"column:exit_point_id"`
	Timestamp     time.Time `json:"timestamp" gorm:"column:timestamp"`
	Status        string    `json:"status" gorm:"column:status"` // e.g., "pending", "approved", "denied"
	ResponseToken string    // Optional: Unique token for secure response validation
}
