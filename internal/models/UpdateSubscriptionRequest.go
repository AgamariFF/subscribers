package models

type UpdateSubscriptionRequest struct {
	// Название сервиса
	ServiceName *string `json:"service_name,omitempty"`
	// Цена подписки
	Price *int `json:"price,omitempty"`
	// Дата начала подписки (формат: 01-2006)
	StartDate *string `json:"start_date,omitempty"`
	// Дата окончания подписки (формат: 01-2006)
	EndDate *string `json:"end_date,omitempty"`
}
