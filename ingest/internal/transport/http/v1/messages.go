package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
	"github.com/AFK068/skilloom/ingest/internal/transport/http/respond"
)

func (h *MessagesHandler) create(c *gin.Context) {
	var req createMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.Messages.SaveMessages(c.Request.Context(), []message.Message{req.toDomain()}); err != nil {
		h.log.Error("save messages failed", zap.Error(err))
		respond.Error(c, http.StatusInternalServerError, "failed to save messages")
		return
	}

	c.Status(http.StatusAccepted)
}
