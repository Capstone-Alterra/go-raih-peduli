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

func (mr *midtransRequest) CreateTransactionBank(IDTransaction string, PaymentType string, Amount int64) (string, string, error) {
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
			return "", "", err
		}

		var vaAccount, validUntil string
		for _, va := range chargeResp.VaNumbers {
			if va.Bank == bank {
				vaAccount = va.VANumber
				validUntil = chargeResp.ExpiryTime
				break
			}
		}

		return vaAccount, validUntil, nil
	}

	return "", "", errors.New("unsupported payment type")
}

func (mr *midtransRequest) CreateTransactionGopay(IDTransaction string, PaymentType string, Amount int64) (string, string, error) {
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
		return "", "", err
	}

	var callback_url, validUntil string
	if len(chargeResp.Actions) > 0 {
		for _, action := range chargeResp.Actions {
			if action.Name == "deeplink-redirect" {
				deepLinkURL := action.URL
				callback_url = deepLinkURL
				break
			}
			validUntil = chargeResp.ExpiryTime
		}
	}

	return callback_url, validUntil, nil
}

func (mr *midtransRequest) CreateTransactionQris(IDTransaction string, PaymentType string, Amount int64) (string, string, error) {
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
		return "", "", err
	}

	var callback_url, validUntil string

	if len(chargeResp.Actions) > 0 {
		for _, action := range chargeResp.Actions {
			if action.Name == "generate-qr-code" {
				callback_url = action.URL
				break
			}
		}
		validUntil = chargeResp.ExpiryTime
	}

	return callback_url, validUntil, nil
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

func (mr *midtransRequest) CheckTransactionStatus(IDTransaction string) (string, error) {
	transactionStatusResp, err := coreapi.CheckTransaction(IDTransaction)
	if err != nil {
		return "", err
	}

	if transactionStatusResp != nil {
		status := mr.TransactionStatus(transactionStatusResp)
		if err != nil {
			return "", err
		}
		return status.Order, nil
	} else {
		return "", nil
	}
}

func (mr *midtransRequest) MappingPaymentName(paymentType string) string {
	switch paymentType {
	case "4":
		return "Bank Permata"
	case "5":
		return "Bank CIMB"
	case "6":
		return "Bank BCA"
	case "7":
		return "Bank BRI"
	case "8":
		return "Bank BNI"
	case "10":
		return "Gopay"
	case "11":
		return "Qris"
	default:
		return "Other"
	}
}
