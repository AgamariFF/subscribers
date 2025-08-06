package services

import (
	"fmt"
	"subscribers/internal/models"
	"subscribers/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func UpdateSubscription(db *gorm.DB, id uuid.UUID, req models.UpdateSubscriptionRequest) error {
	var sub models.Subscription

	if err := db.First(&sub, "id = ?", id).Error; err != nil {
		return err
	}

	if req.ServiceName != nil {
		sub.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		sub.MonthlyPrice = *req.Price
	}
	if req.StartDate != nil {
		startYM, err := utils.ParseYearMonth(*req.StartDate)
		if err != nil {
			return fmt.Errorf("invalid start_date: %w", err)
		}
		sub.StartedAt = startYM
	}

	if req.EndDate != nil {
		if *req.EndDate == "" {
			sub.EndedAt = nil
		} else {
			endYM, err := utils.ParseYearMonth(*req.EndDate)
			if err != nil {
				return fmt.Errorf("invalid end_date: %w", err)
			}
			sub.EndedAt = &endYM
		}
	}

	if err := db.Save(&sub).Error; err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}
