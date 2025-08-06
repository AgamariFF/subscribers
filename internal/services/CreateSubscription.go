package services

import (
	"fmt"
	"subscribers/internal/models"
	"subscribers/internal/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CreateSubscription(db *gorm.DB, req models.CreateSubscriptionRequest) (uuid.UUID, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Invalid UUID: %w", err)
	}

	startYM, err := utils.ParseYearMonth(req.StartDate)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error parsing start_date: %w", err)
	}

	var endYM *models.YearMonth
	if req.EndDate != nil {
		ym, err := utils.ParseYearMonth(*req.EndDate)
		if err != nil {
			return uuid.Nil, fmt.Errorf("Error parsing end_date: %w", err)
		}
		endYM = &ym
	}

	exists, err := utils.SubscriptionExists(db, userID, req.ServiceName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Record search error: %w", err)
	}
	if exists {
		return uuid.Nil, fmt.Errorf("A subscription already exists for this user and the service.")
	}

	sub := models.Subscription{
		ID:           uuid.New(),
		ServiceName:  req.ServiceName,
		MonthlyPrice: req.Price,
		UserID:       userID,
		StartedAt:    startYM,
		EndedAt:      endYM,
	}

	if err := db.Create(&sub).Error; err != nil {
		return uuid.Nil, fmt.Errorf("error saving the subscription: %w", err)
	}

	return sub.ID, nil
}
