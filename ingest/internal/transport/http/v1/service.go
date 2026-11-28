package v1

import (
	"context"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/messages_mock.go -package=mock

type MessagesService interface {
	SaveMessages(ctx context.Context, p principal.Principal, msgs []message.Message) error
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
