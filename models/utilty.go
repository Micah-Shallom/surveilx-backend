package models

import (
	"errors"

	"gorm.io/gorm"
)

func CheckExists(db *gorm.DB, receiver any, query any, args ...any) bool {

	tx := db.Where(query, args...).First(receiver)
	return !errors.Is(tx.Error, gorm.ErrRecordNotFound)
}
