package di

import (
	"database/sql"

	"github.com/ghulammuzz/backend-parkerin/internal/applicants/handler"
	appRepo "github.com/ghulammuzz/backend-parkerin/internal/applicants/repo"
	appSvc "github.com/ghulammuzz/backend-parkerin/internal/applicants/svc"
	storeRepo "github.com/ghulammuzz/backend-parkerin/internal/store/repo"
	storeSvc "github.com/ghulammuzz/backend-parkerin/internal/store/svc"
	userRepo "github.com/ghulammuzz/backend-parkerin/internal/users/repo"
	"github.com/google/wire"
)

func InitializedApplicationServiceFake(sb *sql.DB) *handler.ApplicationHandler {
	wire.Build(
		handler.NewApplicationHandler,
		appSvc.NewApplicationService,
		appRepo.NewApplicationRepository,
		storeSvc.NewStoreService,
		storeRepo.NewStoreRepository,
		userRepo.NewUserRepository,
	)

	return &handler.ApplicationHandler{}
}
