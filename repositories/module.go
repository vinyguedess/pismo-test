package repositories

import (
	"go.uber.org/fx"

	"pismo/providers"
)

var Module = fx.Provide(
	providers.NewDBDialector,
	providers.NewDBConnection,
)
