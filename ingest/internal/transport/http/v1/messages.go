package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
)

func (h *MessagesHandler) create(c *gin.Context) {
	var req createMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := h.Messages.SaveMessages(c.Request.Context(), []message.Message{req.toDomain()}); err != nil {
		h.log.Error("save messages failed", zap.Error(err))
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusAccepted)
}
