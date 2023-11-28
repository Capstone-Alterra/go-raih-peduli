package helpers

import (
	"errors"
	"raihpeduli/config"

	"raihpeduli/features/transaction"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type midtransRequest struct{}

func NewMidtransRequest() MidtransInterface {
	return &midtransRequest{}
}

func (mr *midtransRequest) CreateTransactionBank(IDTransaction string, PaymentType string, Amount int64) (string, error) {
	mtconfig := config.LoadMidtransConfig()

	midtrans.ServerKey = mtconfig.MT_SERVER_KEY
	midtrans.ClientKey = mtconfig.MT_CLIENT_KEY
	midtrans.Environment = midtrans.Sandbox

	bankMap := map[string]string{
		"4": "permata",
		"5": "cimb",
		"6": "bca",
		"7": "bri",
		"8": "bni",
		"9": "mandiri",
	}

	if bank, ok := bankMap[PaymentType]; ok {
		var midtransBank midtrans.Bank

		switch bank {
		case "bca":
			midtransBank = midtrans.BankBca
		case "bri":
			midtransBank = midtrans.BankBri
		case "bni":
			midtransBank = midtrans.BankBni
		case "mandiri":
			midtransBank = midtrans.BankMandiri
		case "cimb":
			midtransBank = midtrans.BankCimb
		case "permata":
			midtransBank = midtrans.BankPermata
		default:
			midtransBank = midtrans.BankBca
		}

		chargeReq := &coreapi.ChargeReq{
			PaymentType:  "bank_transfer",
			BankTransfer: &coreapi.BankTransferDetails{Bank: midtransBank},
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  IDTransaction,
				GrossAmt: Amount,
			},
		}

		chargeResp, err := coreapi.ChargeTransaction(chargeReq)
		if err != nil {
			return "", err
		}

		var vaAccount string
		for _, va := range chargeResp.VaNumbers {
			if va.Bank == bank {
				vaAccount = va.VANumber
				break
			}
		}

		return vaAccount, nil
	}

	return "", errors.New("unsupported payment type")

}

func (mr *midtransRequest) CreateTransactionGopay(IDTransaction string, PaymentType string, Amount int64) (string, error) {
	mtconfig := config.LoadMidtransConfig()

	midtrans.ServerKey = mtconfig.MT_SERVER_KEY
	midtrans.ClientKey = mtconfig.MT_CLIENT_KEY
	midtrans.Environment = midtrans.Sandbox

	chargeReq := &coreapi.ChargeReq{
		Gopay: &coreapi.GopayDetails{
			EnableCallback: true,
		},
		PaymentType: "gopay",
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  IDTransaction,
			GrossAmt: Amount,
		},
	}

	chargeResp, err := coreapi.ChargeTransaction(chargeReq)
	if err != nil {
		return "", err
	}

	var callback_url = ""
	if len(chargeResp.Actions) > 0 {
		for _, action := range chargeResp.Actions {
			if action.Name == "deeplink-redirect" {
				deepLinkURL := action.URL
				callback_url = deepLinkURL
				break
			}
		}
	}

	return callback_url, nil
}

func (mr *midtransRequest) CreateTransactionQris(IDTransaction string, PaymentType string, Amount int64) (string, error) {
	mtconfig := config.LoadMidtransConfig()

	midtrans.ServerKey = mtconfig.MT_SERVER_KEY
	midtrans.ClientKey = mtconfig.MT_CLIENT_KEY
	midtrans.Environment = midtrans.Sandbox

	chargeReq := &coreapi.ChargeReq{
		Qris: &coreapi.QrisDetails{
			Acquirer: "gopay",
		},
		PaymentType: "qris",
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  IDTransaction,
			GrossAmt: Amount,
		},
	}

	chargeResp, err := coreapi.ChargeTransaction(chargeReq)
	if err != nil {
		return "", err
	}

	var url string
	if len(chargeResp.Actions) > 0 {
		for _, action := range chargeResp.Actions {
			if action.Name == "generate-qr-code" {
				url = action.URL
				break
			}
		}
	}

	return url, nil
}

func (mr *midtransRequest) TransactionStatus(transactionStatusResp *coreapi.TransactionStatusResponse) transaction.Status {
	var status transaction.Status

	if transactionStatusResp.TransactionStatus == "capture" {
		if transactionStatusResp.FraudStatus == "challenge" {
			// TODO set transaction status on your database to 'challenge'
			// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			status.Transaction = "challenge"
			status.Order = "2"
		} else if transactionStatusResp.FraudStatus == "accept" {
			// TODO set transaction status on your database to 'success'
			status.Transaction = "success"
			status.Order = "5"
		}
	} else if transactionStatusResp.TransactionStatus == "settlement" {
		// TODO set transaction status on your databaase to 'success'
		status.Transaction = "success"
		status.Order = "5"
	} else if transactionStatusResp.TransactionStatus == "deny" {
		// TODO you can ignore 'deny', because most of the time it allows payment retries
		// and later can become success
		status.Transaction = "success"
		status.Order = "4"
	} else if transactionStatusResp.TransactionStatus == "cancel" || transactionStatusResp.TransactionStatus == "expire" {
		// TODO set transaction status on your databaase to 'failure'
		status.Transaction = "failed"
		status.Order = "3"
	} else if transactionStatusResp.TransactionStatus == "pending" {
		// TODO set transaction status on your databaase to 'pending' / waiting payment
		status.Transaction = "pending"
		status.Order = "2"
	}

	return status
}
