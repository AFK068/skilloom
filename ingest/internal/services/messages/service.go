package messages

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
	"github.com/AFK068/skilloom/ingest/internal/domain/principal"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/publisher_mock.go -package=mock

type Publisher interface {
	PublishMessage(ctx context.Context, p principal.Principal, msg message.Message) error
}

type Deps struct {
	fx.In

	Publisher Publisher
	Log       *zap.Logger
}

type Service struct {
	Deps
	log *zap.Logger
}

func New(d Deps) *Service {
	return &Service{
		Deps: d,
		log:  d.Log.Named("messages service"),
	}
}
