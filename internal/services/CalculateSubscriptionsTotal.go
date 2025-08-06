package services

import (
	"subscribers/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func CalculateSubscriptionsTotal(db *gorm.DB, userID uuid.UUID, serviceName string, startYM, endYM *models.YearMonth) (int, error) {
	var subscriptions []models.Subscription

	query := db.Where("user_id = ?", userID)

	if serviceName != "" {
		query = query.Where("service_name = ?", serviceName)
	}

	if startYM != nil {
		query = query.Where(
			"(EXTRACT(YEAR FROM started_at) * 100 + EXTRACT(MONTH FROM started_at)) >= ?",
			startYM.Year*100+int(startYM.Month),
		)
	}

	if endYM != nil {
		query = query.Where(
			"(EXTRACT(YEAR FROM started_at) * 100 + EXTRACT(MONTH FROM started_at)) <= ?",
			endYM.Year*100+int(endYM.Month),
		)
	}

	if err := query.Find(&subscriptions).Error; err != nil {
		return 0, err
	}

	total := 0
	for _, sub := range subscriptions {
		total += sub.MonthlyPrice
	}

	return total, nil
}
