package handlers

import (
	"go.uber.org/fx"

	"pismo/services"
)

var Module = fx.Provide(
	services.NewAccountService,
	services.NewTransactionService,
)
