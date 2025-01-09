package svc

import (
	"fmt"
	"time"

	paymentEntity "github.com/ghulammuzz/backend-parkerin/internal/payment/entity"
	userRepo "github.com/ghulammuzz/backend-parkerin/internal/users/repo"
	"github.com/ghulammuzz/backend-parkerin/pkg/log"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type PaymentService interface {
	CreateTransaction(userID, packageID int) (*paymentEntity.CreateTransactionResponse, error)
}

type paymentService struct {
	userRepo       userRepo.UserRepository
	midtransClient *snap.Client
}

func (s *paymentService) CreateTransaction(userID, packageID int) (*paymentEntity.CreateTransactionResponse, error) {

	log.Debug("init svc")
	var currentAmount int

	if packageID == 1 {
		currentAmount = 3000
	} else if packageID == 2 {
		currentAmount = 50000
	} else {
		currentAmount = 50000
	}

	log.Debug("current Ammount : ", currentAmount)

	users, err := s.userRepo.Detail(userID)
	if err != nil {
		log.Debug("error user repo detail")
		return nil, err
	}
	custAddress := &midtrans.CustomerAddress{
		FName:       users.Name,
		LName:       "",
		Phone:       users.PhoneNumber,
		Address:     "",
		City:        "",
		Postcode:    "",
		CountryCode: "IDN",
	}

	log.Debug("user name : ", users.Name)

	chargeReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  fmt.Sprintf("ORD-%s-%d", users.Name, time.Now().UnixMilli()),
			GrossAmt: int64(currentAmount),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName:    users.Name,
			BillAddr: custAddress,
			ShipAddr: custAddress,
		},
		EnabledPayments: snap.AllSnapPaymentType,
		Items: &[]midtrans.ItemDetails{
			{
				ID:    "T-Harian",
				Qty:   1,
				Price: int64(currentAmount),
				Name:  "Paket Harian",
			},
		},
	}

	log.Debug("charge req : ", chargeReq)

	chargeRes, errResp := s.midtransClient.CreateTransaction(chargeReq)
	if errResp != nil {
		log.Debug("error charge transaction")
		return nil, err
	}

	transaction := paymentEntity.CreateTransactionResponse{
		Token: chargeRes.Token,
		URL:   chargeRes.RedirectURL,
	}

	return &transaction, nil
}

func NewPaymentService(userRepo userRepo.UserRepository, midtransClient *snap.Client) PaymentService {
	return &paymentService{userRepo: userRepo, midtransClient: midtransClient}
}


/*

get user id by token
get user data by user id
get product data by product id
attach price, product name, id to midtrans-item
attach name, phone to midtrans-customer


*/