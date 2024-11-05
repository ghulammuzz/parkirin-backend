package di

import (
	"database/sql"

	repoApp "github.com/ghulammuzz/backend-parkerin/internal/applicants/repo"
	"github.com/ghulammuzz/backend-parkerin/internal/store/handler"
	repoStore "github.com/ghulammuzz/backend-parkerin/internal/store/repo"
	"github.com/ghulammuzz/backend-parkerin/internal/store/svc"
	repoUser "github.com/ghulammuzz/backend-parkerin/internal/users/repo"
	"github.com/google/wire"
)

func InitializedStoreServiceFake(sb *sql.DB) *handler.StoreHandler {
	wire.Build(
		handler.NewStoreHandler,
		svc.NewStoreService,
		repoStore.NewStoreRepository,
		repoUser.NewUserRepository,
		repoApp.NewApplicationRepository,
	)

	return &handler.StoreHandler{}
}
