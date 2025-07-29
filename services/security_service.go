package services

import (
	"survielx-backend/database"
	"survielx-backend/models"
)

func CreateSecurity(security *models.Security) error {
	if err := database.DB.Create(security).Error; err != nil {
		return err
	}
	return nil
}
