package utils

import (
	"subscribers/internal/models"
	"time"
)

func ParseYearMonth(s string) (models.YearMonth, error) {
	t, err := time.Parse("01-2006", s)
	if err != nil {
		t, err = time.Parse("2006-01", s)
	}
	if err != nil {
		return models.YearMonth{}, err
	}
	return models.YearMonth{Year: t.Year(), Month: t.Month()}, nil
}
