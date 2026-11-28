package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
	"github.com/AFK068/skilloom/ingest/internal/services/messages"
	"github.com/AFK068/skilloom/ingest/internal/transport/http/middleware"
	"github.com/AFK068/skilloom/ingest/internal/transport/http/respond"
)

func (h *MessagesHandler) create(c *gin.Context) {
	p, ok := middleware.Principal(c)
	if !ok {
		respond.Error(c, http.StatusUnauthorized, "missing principal")
		return
	}

	var req createMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respond.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	err := h.Messages.SaveMessages(c.Request.Context(), p, []message.Message{req.toDomain()})
	switch {
	case errors.Is(err, messages.ErrInvalidMessage):
		respond.Error(c, http.StatusBadRequest, "invalid message")
		return
	case err != nil:
		h.log.Error("save messages failed", zap.Error(err))
		respond.Error(c, http.StatusInternalServerError, "failed to save messages")
		return
	}

	c.Status(http.StatusAccepted)
}
