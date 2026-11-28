package messages

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/AFK068/skilloom/ingest/internal/domain/message"
)

type RedisStore interface {
	Get(ctx context.Context, key string) (string, error)
}

type Publisher interface {
	Publish(ctx context.Context, subject, msgID string, data []byte) error
}

type Deps struct {
	fx.In

	Redis     RedisStore
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

func (s *Service) SaveMessages(_ context.Context, _ []message.Message) error {
	s.log.Info("save messages stub")
	return nil
}
