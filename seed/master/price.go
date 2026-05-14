package master

import (
	"backend-brevet/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedPrices is a function that seeds prices to the database
func SeedPrices(db *gorm.DB) error {

	prices := []models.Price{
		{ID: uuid.New(), GroupType: models.MahasiswaGunadarma, Price: 750000},
		{ID: uuid.New(), GroupType: models.MahasiswaNonGunadarma, Price: 1000000},
		{ID: uuid.New(), GroupType: models.Umum, Price: 2300000},
	}

	for _, price := range prices {
		var existing models.Price
		err := db.Where("group_type = ?", price.GroupType).First(&existing).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := db.Create(&price).Error; err != nil {
				return err
			}
			continue
		}
		if err != nil {
			return err
		}
		if err := db.Model(&existing).Update("price", price.Price).Error; err != nil {
			return err
		}
	}
	return nil
}
