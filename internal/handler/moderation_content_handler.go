package handler

import (
	// "golang_gpt/internal/entity"
	"golang_gpt/internal/entity"
	"golang_gpt/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

type ModerationHandler struct {
	service *service.ContentService
	server  *socketio.Server
}

func NewModerationHandler(service *service.ContentService, server *socketio.Server) *ModerationHandler {
	return &ModerationHandler{service: service, server: server}
}


// GetContentForModeration возвращает список контента, ожидающего модерации
func (h *ModerationHandler) GetContentForModeration(c *gin.Context) {
	contents, err := h.service.GetContentForModeration()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка получения данных"})
		return
	}
	log.Printf("contents: %v", contents)
	c.JSON(http.StatusOK, contents)
}


// ModerateContent позволяет модератору принять или отклонить контент
func (h *ModerationHandler) ModerateContent(c *gin.Context) {
	var req entity.ModerateContentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный ввод"})
		return
	}

	// Обновляем статус контента в истории
	err := h.service.UpdateContentLatestHistory(req.ContentID, req.StatusID)
	if err != nil {
		log.Printf("Ошибка обновления статуса контента ID %d: %v", req.ContentID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка обновления статуса", "details": err.Error()})
		return
	}

	// Если контент принят, отправляем его на монитор
	if req.StatusID == entity.ContentApproved {
		content, _ := h.service.GetContentByID(req.ContentID)
		approvedContent := &entity.ContentForMonitor{
			FileName:  content.FileName,
			FilePath:  content.FilePath,
			StartTime: content.StartTime,
			EndTime:   content.EndTime,
		}

		// Отправляем данные на устройство
		h.server.BroadcastToRoom("/", content.MacAddress, "data", approvedContent)
	}

	// Если контент отклонен, удаляем его
	// if req.StatusID == entity.ContentRejected {
	// 	h.service.DeleteContent(req.ContentID)
	// }

	c.JSON(http.StatusOK, gin.H{"message": "Модерация завершена"})
}
