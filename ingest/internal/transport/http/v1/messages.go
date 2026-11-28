package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MessagesDeps struct {
	fx.In

	Log *zap.Logger
}

type MessagesHandler struct {
	MessagesDeps
	log *zap.Logger
}

func NewMessagesHandler(d MessagesDeps) *MessagesHandler {
	return &MessagesHandler{
		MessagesDeps: d,
		log:          d.Log.Named("messages http handler"),
	}
}

func (h *MessagesHandler) Register(rg *gin.RouterGroup) {
	rg.POST("/messages", h.create)
}

func (h *MessagesHandler) create(c *gin.Context) {
	h.log.Info("message stub received")
	c.Status(http.StatusAccepted)
}
