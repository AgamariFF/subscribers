package services

import (
	"subscribers/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetSubscriptions(db *gorm.DB, userID uuid.UUID) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	if err := db.Where("user_id = ?", userID).Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}
