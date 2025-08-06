package services

import (
	"subscribers/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GetSubscriptionByID(db *gorm.DB, id uuid.UUID) (*models.Subscription, error) {
	var sub models.Subscription
	err := db.First(&sub, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}
