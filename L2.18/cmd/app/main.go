package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"
	"time"

	"task-manager/internal/config"
	"task-manager/internal/handler"
	"task-manager/internal/logger"
	"task-manager/internal/repo"
	"task-manager/internal/service"
)

func main() {
	// Загружаем конфиг
	cfg := config.Load()
	log.Printf("Starting server on port %s...", cfg.AppPort)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Подключение к БД
	db, err := repo.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(ctx); err != nil {
		log.Fatalf("Database not reachable: %v", err)
	}
	log.Println("Connected to database successfully")

	// Инициализация сервиса
	taskService := service.NewService(db)
	// Получаем middleware для логирования
	requestLoggerMiddleware := logger.RequestLogger()
	// Инициализация handler
	ginHandler := handler.NewHandler(taskService, requestLoggerMiddleware)

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Запуск HTTP сервера
	if err := ginHandler.Run(shutdownCtx, cfg.AppPort); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
