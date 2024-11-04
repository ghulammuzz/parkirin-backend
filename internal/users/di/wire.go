package di

import (
	"database/sql"

	"github.com/ghulammuzz/backend-parkerin/internal/users/handler"
	"github.com/ghulammuzz/backend-parkerin/internal/users/repo"
	"github.com/ghulammuzz/backend-parkerin/internal/users/svc"
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
)

func InitializedUsersServiceFake(sb *sql.DB, val *validator.Validate) *handler.UserHandler {
	wire.Build(
		handler.NewUserHandler,
		svc.NewUserService,
		repo.NewUserRepository,
	)

	return &handler.UserHandler{}
}
