package services

import (
	"subscribers/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func DeleteSubscription(db *gorm.DB, id uuid.UUID) error {
	result := db.Delete(&models.Subscription{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
