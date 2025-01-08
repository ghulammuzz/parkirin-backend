package config

import (
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

func InitMidtrans() *snap.Client {
	midtrans.ServerKey = os.Getenv("MIDTRANS_SERVER_KEY")
	// midtrans.ClientKey = os.Getenv("MIDTRANS_CLIENT_KEY")

	midtrans.Environment = midtrans.Sandbox

	snapClient := snap.Client{}
	snapClient.New(midtrans.ServerKey, midtrans.Sandbox)

	return &snapClient
}
