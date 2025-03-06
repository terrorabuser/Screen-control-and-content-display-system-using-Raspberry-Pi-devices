package main

import (
	"log"

	"golang_gpt/internal/app"
)

func main() {
	// Инициализируем приложение
	application, err := app.NewApp()
	if err != nil {
		log.Fatal("Ошибка инициализации приложения:", err)
	}

	// Запускаем сервер
	application.Router.Run(":8080")
}
