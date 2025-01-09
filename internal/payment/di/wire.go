package di

import (
	"database/sql"

	"github.com/ghulammuzz/backend-parkerin/internal/payment/handler"
	paySvc "github.com/ghulammuzz/backend-parkerin/internal/payment/svc"
	userRepo "github.com/ghulammuzz/backend-parkerin/internal/users/repo"
	"github.com/google/wire"
	"github.com/midtrans/midtrans-go/snap"
)

func InitializedPaymentServiceFake(sb *sql.DB, midtransClient *snap.Client) *handler.PaymentHandler {
	wire.Build(
		handler.NewPaymentHandler,
		paySvc.NewPaymentService,
		userRepo.NewUserRepository,
	)

	return &handler.PaymentHandler{}
}
