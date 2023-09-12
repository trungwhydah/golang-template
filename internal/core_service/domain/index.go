package domain

import (
	"github.com/golang/be/internal/core_service/domain/product"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(product.NewUseCase),
)
