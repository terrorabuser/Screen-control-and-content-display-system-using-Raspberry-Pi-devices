package app

import (
	"database/sql"
	"fmt"

	"golang_gpt/internal/config"
	"golang_gpt/internal/handler"
	"golang_gpt/internal/repository"
	"golang_gpt/internal/server"
	"golang_gpt/internal/service"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/googollee/go-socket.io"
)

// App содержит все зависимости приложения
type App struct {
    DB           *sql.DB
    Router       *gin.Engine
    SocketServer *socketio.Server
}

// NewApp инициализирует приложение
func NewApp() (*App, error) {
	// Загружаем конфиг
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфигурации: %w", err)
	}

	// Подключение к БД
	db, err := sql.Open("postgres", fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.Password, cfg.Database.DBName, cfg.Database.SSLMode,
	))
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// Инициализация репозиториев
	monitorRepo := repository.NewMonitorRepository(db)
	contentRepo := repository.NewContentRepository(db)
	apiRepo := repository.NewApiRepository(db)
	socketRepo := repository.NewSocketRepository(db)

	// Инициализация сервисов
	authMonitorService := service.NewAuthMonitorService(monitorRepo)
	contentService := service.NewContentService(contentRepo)

	// Инициализация обработчиков
	monitorHandler := handler.NewMonitorHandler(monitorRepo)
	authMonitorHandler := handler.NewAuthMonitorHandler(authMonitorService)
	socketHandler := handler.NewSocketHandler(socketRepo, authMonitorService)

	// Инициализация Socket.IO сервера
	socketServer := server.NewSocketServer(socketHandler)

	// Обработчики с сокетом
	moderationHandler := handler.NewModerationHandler(contentService, socketServer)
	contentHandler := handler.NewContentHandler(contentService, socketServer)
	apiHandler := handler.NewApiHandler(apiRepo)

	// Настраиваем маршрутизатор
	router := server.SetupRouter(monitorHandler, contentHandler, moderationHandler, socketServer, authMonitorHandler, apiHandler)

	return &App{
		DB:           db,
		Router:       router,
		SocketServer: socketServer,
	}, nil
}
