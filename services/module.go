package services

import (
	"go.uber.org/fx"

	"pismo/repositories"
)

var Module = fx.Provide(
	repositories.NewAccountRepository,
	repositories.NewOperationTypeRepository,
	repositories.NewTransactionRepository,
)
