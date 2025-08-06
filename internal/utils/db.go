package utils

import (
	"errors"
	"subscribers/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func SubscriptionExists(db *gorm.DB, userID uuid.UUID, serviceName string) (bool, error) {
	var existing models.Subscription
	err := db.Where("user_id = ? AND service_name = ?", userID, serviceName).First(&existing).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, err
}
