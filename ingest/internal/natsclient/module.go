package natsclient

import "go.uber.org/fx"

var Module = fx.Module(
	"nats",

	fx.Provide(
		New,
	),
	fx.Invoke(
		func(*Client) {},
	),
)
