package models

// CreateSubscriptionRequest represents запрос на создание подписки
// swagger:model CreateSubscriptionRequest
type CreateSubscriptionRequest struct {
	// Название сервиса
	// required: true
	ServiceName string `json:"service_name" binding:"required"`
	// Цена подписки (мин 0)
	// required: true
	Price int `json:"price" binding:"required,min=0"`
	// ID пользователя в формате UUID
	// required: true
	UserID string `json:"user_id" binding:"required,uuid"`
	// Дата начала подписки (формат: 2006-01 или 01-2006)
	// required: true
	StartDate string `json:"start_date" binding:"required"`
	// Дата окончания подписки (формат: 2006-01 или 01-2006)
	// required: false
	EndDate *string `json:"end_date,omitempty"`
}
