package handler

import (
	"golang_gpt/internal/entity"
	"golang_gpt/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type ContentHandler struct {
	service *service.ContentService
	server *socketio.Server
}

func NewContentHandler(service *service.ContentService, server *socketio.Server) *ContentHandler {
    return &ContentHandler{service: service, server: server}
}

// Добавление контента
func (h *ContentHandler) AddContent(c *gin.Context) {
	
	req := entity.ContentRequest{
		Building:  c.PostForm("building"),
		Floor:     c.PostForm("floor"),
		Notes:     c.PostForm("notes"),
		FileName:  c.PostForm("file_name"),
		FilePath:  c.PostForm("file_path"),
		StartTime: c.PostForm("start_time"),
		EndTime:   c.PostForm("end_time"),
	}
	 // 2. Получаем загруженный файл
	file, err := c.FormFile("file")
	if err != nil {
		 c.JSON(http.StatusBadRequest, gin.H{"error": "File not found"})
		 return
	}

    // 3. Генерируем путь для сохранения файла
    filePath := "./uploads/" + file.Filename

    // 4. Сохраняем файл
    if err := c.SaveUploadedFile(file, filePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
        return
    }

	// Получаем ID пользователя из контекста

	// userId, exits := c.Get(userId)
	// log.Printf("userId: %s", userId)
	// if !exits {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	// 	return
	// }
	validUserID := int64(1)
	// validUserID, ok := userId.(int64)
	// if !ok {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "unauthorized"})
	// 	return
	// }
	


	
	// Получаем MAC-адрес по локации
	MacAddress, err := h.service.GetMacAddressByLocation(req.Building, req.Floor, req.Notes)
	if err != nil {
		log.Println("Ошибка получения MAC-адреса:", err)
		return
	}

	// Преобразуем в entity.Content
	content := &entity.ContentForDB{
		UserID:  validUserID,
		MacAddress:      MacAddress,
		FileName:  req.FileName,
		FilePath:  req.FilePath,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Добавляем контент
	id, err := h.service.AddContent(content)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Content sended to moderation", "id": id})
}
