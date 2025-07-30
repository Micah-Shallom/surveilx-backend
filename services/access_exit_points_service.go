package services

import (
	"survielx-backend/database"
	"survielx-backend/models"
)

func CreateAccessExitPoint(point *models.AccessExitPoint) error {
	result := database.DB.Create(point)
	return result.Error
}

func GetAccessExitPoints(points *[]models.AccessExitPoint) error {
	result := database.DB.Find(points)
	return result.Error
}

func GetAccessExitPoint(id string, point *models.AccessExitPoint) error {
	result := database.DB.First(point, "id = ?", id)
	return result.Error
}

func UpdateAccessExitPoint(point *models.AccessExitPoint) error {
	result := database.DB.Save(point)
	return result.Error
}

func DeleteAccessExitPoint(id string) error {
	result := database.DB.Delete(&models.AccessExitPoint{}, "id = ?", id)
	return result.Error
}
