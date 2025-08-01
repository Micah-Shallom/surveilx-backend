package services

import (
	"survielx-backend/models"

	"gorm.io/gorm"
)

func CreateAccessExitPoint(db *gorm.DB, point *models.AccessExitPoint) error {
	result := db.Create(point)
	return result.Error
}

func GetAccessExitPoints(db *gorm.DB, points *[]models.AccessExitPoint) error {
	result := db.Find(points)
	return result.Error
}

func GetAccessExitPoint(db *gorm.DB, id string, point *models.AccessExitPoint) error {
	result := db.First(point, "id = ?", id)
	return result.Error
}

func UpdateAccessExitPoint(db *gorm.DB, point *models.AccessExitPoint) error {
	result := db.Save(point)
	return result.Error
}

func DeleteAccessExitPoint(db *gorm.DB, id string) error {
	result := db.Delete(&models.AccessExitPoint{}, "id = ?", id)
	return result.Error
}
