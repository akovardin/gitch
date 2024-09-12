package tasks

import "go.uber.org/fx"

var Module = fx.Module(
	"tasks",
	fx.Provide(NewSync),
)
