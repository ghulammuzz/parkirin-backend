package entity

import "time"

type Transaction struct {
	ID              int
	UserID          int
	PackageID       int
	Status          string
	Amount          int
	TransactionTime time.Time
	PaymentURL      string
	MidtransID      string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Payment struct {
	ID            int
	TransactionID int
	PaymentMethod string
	PaymentStatus string
	PaymentTime   time.Time
	RawResponse   string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type CreateTransactionResponse struct {
	Token string
	URL   string
}
