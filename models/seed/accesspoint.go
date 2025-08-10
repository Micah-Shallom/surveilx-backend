package seed

import (
	"fmt"
	"survielx-backend/models"
	"survielx-backend/utility"

	"gorm.io/gorm"
)

func SeedAccessPoint(db *gorm.DB) {
	var (
		ap    models.AccessExitPoint
		count int64
	)

	if err := db.Model(&ap).Where("name IN ?", []string{"North Gate", "Main Gate"}).Count(&count).Error; err != nil {
		fmt.Println("accesspoint seeding: " + err.Error())
		return
	}

	if count > 0 {
		fmt.Println("AccessPoint already exist, skipping seeding...")
		return
	} else {

		accessexitpoint := []models.AccessExitPoint{
			{
				ID:   utility.GenerateUUID(),
				Name: "North Gate",
			},
			{
				ID:   utility.GenerateUUID(),
				Name: "Main Gate",
			},
		}

		db = db.Debug()

		for _, aep := range accessexitpoint {
			if err := db.Create(&aep).Error; err != nil {
				fmt.Println("failed to seed plan: " + err.Error())
			}
		}
	}
}
