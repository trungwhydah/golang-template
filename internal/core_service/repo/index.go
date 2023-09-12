package repo

import (
	"github.com/golang/be/internal/core_service/repo/product"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(product.NewMongoRepo),
)
