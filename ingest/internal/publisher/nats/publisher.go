package natspub

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

//go:generate go tool mockgen -source=$GOFILE -destination=mock/client_mock.go -package=mock

type Client interface {
	Publish(ctx context.Context, subject, msgID string, data []byte) error
}

type Deps struct {
	fx.In

	Client Client
	Log    *zap.Logger
}

type Publisher struct {
	Deps
	log *zap.Logger
}

func New(d Deps) *Publisher {
	return &Publisher{
		Deps: d,
		log:  d.Log.Named("nats publisher"),
	}
}
