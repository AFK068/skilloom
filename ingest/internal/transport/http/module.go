package httpserver

import "go.uber.org/fx"

var Module = fx.Module(
	"http",

	fx.Provide(
		New,
	),
	fx.Invoke(
		(*Server).Start,
	),
)
