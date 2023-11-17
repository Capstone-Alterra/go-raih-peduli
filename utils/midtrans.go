package utils

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"

	"raihpeduli/config"
)

func MidtransCoreAPIClient() coreapi.Client {
	mtconfig := config.LoadMidtransConfig()

	var coreAPIClient coreapi.Client
	coreAPIClient.New(mtconfig.MT_SERVER_KEY, midtrans.Sandbox)

	return coreAPIClient
}
