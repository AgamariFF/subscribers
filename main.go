package main

import (
	"fmt"
	"subscribers/config"
	"subscribers/internal/db"
	"subscribers/internal/handlers"
	"subscribers/logger"

	_ "subscribers/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Subscriptions API
// @version 1.0
// @description API для управления подписками
// @host localhost:8080
func main() {
	cfg := config.LoadConfig()

	logger.InitLogger(cfg.LogLevel)
	defer logger.SugaredLogger.Sync()

	logger.SugaredLogger.Info("The logger is initialized")

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	gormDB := db.ConnectGORM(dsn)
	if gormDB == nil {
		logger.SugaredLogger.Fatal("Couldn't connect to the database, shutting down")
	}

	if err := db.AutoMigrate(gormDB); err != nil {
		logger.SugaredLogger.Errorf("Migration error: %v", err)
	}

	router := gin.Default()

	router.Static("/docs", "./docs")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/createSubscription", handlers.CreateSubscriptionHandler(gormDB))
	router.GET("/subscriptions", handlers.GetSubscriptionsHandler(gormDB))
	router.GET("/subscriptions/:id", handlers.GetSubscriptionHandler(gormDB))
	router.PATCH("/subscriptions/:id", handlers.UpdateSubscriptionHandler(gormDB))
	router.DELETE("/subscriptions/:id", handlers.DeleteSubscriptionHandler(gormDB))
	router.GET("/subscriptions/total", handlers.GetSubscriptionsTotalHandler(gormDB))
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/swagger/index.html")
	})

	logger.SugaredLogger.Infof("The server is running and listening on the port %s", cfg.AppPort)

	if err := router.Run(":" + cfg.AppPort); err != nil {
		logger.SugaredLogger.Fatalf("Failed to start server: %v", err)
	}
}
