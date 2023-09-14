package app

import (
	cmconfig "github.com/golang/be/config/common"
	config "github.com/golang/be/config/core_service"
	cmdomain "github.com/golang/be/internal/common/domain"
	cmrepo "github.com/golang/be/internal/common/repo"
	"github.com/golang/be/internal/core_service/api"
	"github.com/golang/be/internal/core_service/api/handler"
	"github.com/golang/be/internal/core_service/api/middleware"
	"github.com/golang/be/internal/core_service/domain"
	"github.com/golang/be/internal/core_service/repo"
	"github.com/golang/be/pkg/common/logger"
	"github.com/golang/be/pkg/common/mongo"
	"github.com/golang/be/pkg/common/msgtranslate"
	"github.com/golang/be/pkg/core_service/firebase"
	"github.com/golang/be/pkg/core_service/firebase/storage"
	"go.uber.org/fx"
)

var InternalOptions = fx.Options(
	// Common Config
	fx.Provide(cmconfig.NewConfig),

	// Config
	fx.Provide(config.NewConfig),

	// Server
	fx.Provide(NewServer),

	// Router
	fx.Provide(api.NewRouter),

	// Controller
	handler.Module,

	// Middleware
	middleware.Module,

	// Use Case
	domain.Module,

	// Repo
	repo.Module,

	// Common Repo
	cmrepo.Module,

	// Common Domain
	cmdomain.Module,
)

var PackageOptions = fx.Options(
	// Mongo
	fx.Provide(mongo.New),

	// Firebase
	fx.Provide(firebase.NewApps),

	// Storage
	fx.Provide(storage.NewBucketHandler),

	// Logger
	fx.Provide(logger.Init),

	// Msg translate
	fx.Invoke(msgtranslate.Init),
)
