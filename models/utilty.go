package models

import (
	"errors"

	"gorm.io/gorm"
)

func CheckExists(db *gorm.DB, receiver any, query any, args ...any) bool {
	tx := db.Where(query, args...).First(receiver)
	return !errors.Is(tx.Error, gorm.ErrRecordNotFound)
}


func UpdateFields(db *gorm.DB, model any, updates any, query any, args ...any) (*gorm.DB, error) {
	result := db.Model(model).Where(query, args...).Updates(updates)
	if result.Error != nil {
		return result, result.Error
	}
	return result, nil
}