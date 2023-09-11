package user

import depinjection "github.com/golang/be/pkg/common/dep_injection"

var Module = depinjection.BulkProvide(
	[]any{},
	"user-controller",
)
