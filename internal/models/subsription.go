package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Subscription представляет модель подписки
// swagger:model Subscription
// SubscriptionSwagger — структура для отображения подписки в Swagger
type SubscriptionSwagger struct {
	ID           string  `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	ServiceName  string  `json:"service_name" example:"Netflix"`
	MonthlyPrice int     `json:"monthly_price" example:"1000"`
	UserID       string  `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174001"`
	StartedAt    string  `json:"started_at" example:"01-2024"`
	EndedAt      *string `json:"ended_at,omitempty" example:"06-2024"`
}

type Subscription struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	ServiceName  string     `json:"service_name"`
	MonthlyPrice int        `json:"monthly_price"`
	UserID       uuid.UUID  `json:"user_id"`
	StartedAt    YearMonth  `json:"started_at"`
	EndedAt      *YearMonth `json:"ended_at,omitempty"`
}

type YearMonth struct {
	Year  int
	Month time.Month
}

func (ym YearMonth) GormDataType() string {
	return "timestamp"
}

func (ym YearMonth) Value() (driver.Value, error) {
	return time.Date(ym.Year, ym.Month, 1, 0, 0, 0, 0, time.UTC), nil
}

func (ym *YearMonth) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var t time.Time
	switch v := value.(type) {
	case time.Time:
		t = v
	case string:
		var err error
		t, err = time.Parse("2006-01", v)
		if err != nil {
			return fmt.Errorf("Не получилосб распарсить YearMonth из строки: %v", err)
		}
	default:
		return fmt.Errorf("Невозможно извлечь данные типо %T в YearMonth", value)
	}
	ym.Year = t.Year()
	ym.Month = t.Month()
	return nil
}

func (ym YearMonth) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%04d-%02d"`, ym.Year, ym.Month)), nil
}

func (ym *YearMonth) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		return nil
	}
	t, err := time.Parse("2006-01", s)
	if err != nil {
		return fmt.Errorf("не удается распарсить YearMonth из JSON: %v", err)
	}
	ym.Year = t.Year()
	ym.Month = t.Month()
	return nil
}
