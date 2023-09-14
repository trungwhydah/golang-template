package user

import (
	"github.com/golang/be/internal/core_service/api/handler/user/product"
	depinjection "github.com/golang/be/pkg/common/dep_injection"
)

var Module = depinjection.BulkProvide(
	[]any{
		product.NewController,
	},
	"user-controller",
)
