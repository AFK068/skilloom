package v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
)

type MessagesService interface {
	SaveMessages(ctx context.Context, msgs []message.Message) error
}

type MessagesDeps struct {
	fx.In

	Messages MessagesService
	Log      *zap.Logger
}

type MessagesHandler struct {
	MessagesDeps
	log *zap.Logger
}

func NewMessagesHandler(d MessagesDeps) *MessagesHandler {
	return &MessagesHandler{
		MessagesDeps: d,
		log:          d.Log.Named("messages"),
	}
}

func (h *MessagesHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/messages", h.create)
}
