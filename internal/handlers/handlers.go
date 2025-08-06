package handlers

import (
	"errors"
	"net/http"
	"subscribers/internal/models"
	"subscribers/internal/services"
	"subscribers/internal/utils"
	"subscribers/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateSubscription
// @Summary      Создаёт новую подписку
// @Description  Добавляет новую подписку пользователя на сервис. Возвращает 409 при дублировании.
// @Tags         subscription
// @Accept       json
// @Produce      json
// @Param        request body models.CreateSubscriptionRequest true "Данные для создания подписки"
// @Success      201 {object} map[string]string "Подписка успешно создана"
// @Failure      400 {object} map[string]string "Неверные данные запроса"
// @Failure      409 {object} map[string]string "Подписка уже существует"
// @Failure      500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router       /createSubscription [post]
func CreateSubscriptionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.SugaredLogger.Info("Subscription addition started")

		var request models.CreateSubscriptionRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			logger.SugaredLogger.Warnf("Bad request when creating a book: " + err.Error())
			return
		} else {
			logger.SugaredLogger.Info("Successfully decoded request")
		}

		subId, err := services.CreateSubscription(db, request)
		if err != nil {
			logger.SugaredLogger.Warnf("Failed to create subscription: %v", err)

			status := http.StatusBadRequest
			if err.Error() == "A subscription already exists for this user and the service." {
				status = http.StatusConflict
			}
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Подписка успешно создана", "id подписки": subId})
	}
}

// GetSubscriptions
// @Summary      Получить список подписок по user_id
// @Tags         subscription
// @Accept       json
// @Produce      json
// @Param        user_id query string true "ID пользователя (UUID)"
// @Success      200 {array} models.SubscriptionSwagger
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions [get]
func GetSubscriptionsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.SugaredLogger.Info("Get subscriptions started")
		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			logger.SugaredLogger.Warn("Bad request when getting subscriptions, absent user_id")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.SugaredLogger.Warn("Bad request when getting subscriptions, invalid UUID")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
			return
		}

		subscriptions, err := services.GetSubscriptions(db, userID)
		if err != nil {
			logger.SugaredLogger.Errorf("Error getting subscriptions: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch subsriptions"})
			return
		}

		c.JSON(http.StatusOK, subscriptions)
		logger.SugaredLogger.Info("Get subscriptions success")
	}
}

// GetSubscription
// @Summary Получить одну подписку по ID
// @Tags subscription
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Success 200 {object} models.SubscriptionSwagger
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [get]
func GetSubscriptionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.SugaredLogger.Info("Get one subscription by id started")
		idStr := c.Param("id")
		subID, err := uuid.Parse(idStr)
		if err != nil {
			logger.SugaredLogger.Warnf("Invalid subscription UUID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID"})
			return
		}

		sub, err := services.GetSubscriptionByID(db, subID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.SugaredLogger.Warnf("Get one subscription by id failed: %v", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			} else {
				logger.SugaredLogger.Errorf("Error fetching subscription: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch subscription"})
			}
			return
		}

		logger.SugaredLogger.Info("Get one subscription by id success")
		c.JSON(http.StatusOK, sub)
	}
}

// UpdateSubscription
// @Summary Обновить подписку
// @Description Для удаления даты окончания подписки необходимо передать пустую строку в поле `ended_at`.
// @Tags subscription
// @Accept json
// @Produce json
// @Param id path string true "ID подписки (UUID)"
// @Param subscription body models.UpdateSubscriptionRequest true "Данные для обновления"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [patch]
func UpdateSubscriptionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.SugaredLogger.Info("Update subscription started")

		idStr := c.Param("id")
		subID, err := uuid.Parse(idStr)
		if err != nil {
			logger.SugaredLogger.Warnf("Invalid subscription ID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID"})
			return
		}

		var req models.UpdateSubscriptionRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			logger.SugaredLogger.Warnf("Invalid request body: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		if err := services.UpdateSubscription(db, subID, req); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.SugaredLogger.Warnf("Subscription not found: %v", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			} else {
				logger.SugaredLogger.Errorf("error UpdateSubscription: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		logger.SugaredLogger.Info("Update subscription success")
		c.JSON(http.StatusOK, gin.H{"message": "subscription updated successfully"})
	}
}

// DeleteSubscription
// @Summary      Удалить подписку по ID
// @Tags         subscription
// @Accept       json
// @Produce      json
// @Param        id path string true "ID подписки (UUID)"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/{id} [delete]
func DeleteSubscriptionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.SugaredLogger.Info("Delete subscription started")

		idStr := c.Param("id")
		subID, err := uuid.Parse(idStr)
		if err != nil {
			logger.SugaredLogger.Warnf("Invalid subscription ID: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription ID"})
			return
		}

		err = services.DeleteSubscription(db, subID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				logger.SugaredLogger.Warnf("subscription not found: %v", err)
				c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
			} else {
				logger.SugaredLogger.Errorf("Failed to delete subscription: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete subscription"})
			}
			return
		}

		logger.SugaredLogger.Info("Delete subscription success")
		c.JSON(http.StatusOK, gin.H{"message": "subscription deleted successfully"})
	}
}

// GetSubscriptionsTotal
// @Summary      Получить суммарную стоимость подписок за период с фильтрацией
// @Tags         subscription
// @Accept       json
// @Produce      json
// @Param        user_id query string true "ID пользователя (UUID)"
// @Param        service_name query string false "Название подписки (фильтр)"
// @Param        start_date query string false "Начало периода (формат: 01-2006)"
// @Param        end_date query string false "Конец периода (формат: 01-2006)"
// @Success      200 {object} map[string]int "Суммарная стоимость"
// @Failure      400 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /subscriptions/total [get]
func GetSubscriptionsTotalHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.SugaredLogger.Info("Get total subscription started")

		userIDStr := c.Query("user_id")
		if userIDStr == "" {
			logger.SugaredLogger.Warnf("user id is empty")
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
			return
		}

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			logger.SugaredLogger.Warnf("invalid UUID format: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid UUID format"})
			return
		}

		serviceName := c.Query("service_name")
		startDateStr := c.Query("start_date")
		endDateStr := c.Query("end_date")

		var startYM *models.YearMonth
		if startDateStr != "" {
			ym, err := utils.ParseYearMonth(startDateStr)
			if err != nil {
				logger.SugaredLogger.Warnf("invalid start_date format: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
				return
			}
			startYM = &ym
		}

		var endYM *models.YearMonth
		if endDateStr != "" {
			ym, err := utils.ParseYearMonth(endDateStr)
			if err != nil {
				logger.SugaredLogger.Warnf("invalid end_date format: %v", err)
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
				return
			}
			endYM = &ym
		}

		total, err := services.CalculateSubscriptionsTotal(db, userID, serviceName, startYM, endYM)
		if err != nil {
			logger.SugaredLogger.Errorf("Failed to calculate total: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate total"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"total_price": total})
	}
}
