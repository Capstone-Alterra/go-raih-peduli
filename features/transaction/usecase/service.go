package usecase

import (
	"errors"
	"raihpeduli/features/transaction"
	"raihpeduli/features/transaction/dtos"
	"raihpeduli/helpers"
	"strconv"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
	"github.com/midtrans/midtrans-go/coreapi"
)

type service struct {
	model         transaction.Repository
	generator     helpers.GeneratorInterface
	mtRequest     helpers.MidtransInterface
	coreAPIClient coreapi.Client
}

func New(model transaction.Repository, generator helpers.GeneratorInterface, mtRequest helpers.MidtransInterface, coreAPIClient coreapi.Client) transaction.Usecase {
	return &service{
		model:         model,
		generator:     generator,
		mtRequest:     mtRequest,
		coreAPIClient: coreAPIClient,
	}
}

func (svc *service) FindAll(page, size int) []dtos.ResTransaction {
	var transactions []dtos.ResTransaction

	transactionsEnt := svc.model.Paginate(page, size)

	for _, transaction := range transactionsEnt {
		var data dtos.ResTransaction

		if err := smapping.FillStruct(&data, smapping.MapFields(transaction)); err != nil {
			log.Error(err.Error())
		}

		transactions = append(transactions, data)
	}

	return transactions
}

func (svc *service) FindByID(transactionID int) *dtos.ResTransaction {
	res := dtos.ResTransaction{}
	transaction := svc.model.SelectByID(transactionID)

	if transaction == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(transaction))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}
func (svc *service) Create(userID int, newTransaction dtos.InputTransaction) (*dtos.ResTransaction, error) {
	transaction := transaction.Transaction{}
	resTransaction := dtos.ResTransaction{}

	err := smapping.FillStruct(&transaction, smapping.MapFields(newTransaction))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	transaction.UserID = userID
	transaction.ID = svc.generator.GenerateRandomID()
	transaction.Status = "1"
	transactionID := svc.model.Insert(transaction)
	if transactionID == -1 {
		return nil, err
	}

	switch transaction.PaymentType {
	case "6", "7", "8":
		req, err := svc.mtRequest.CreateTransactionBank(strconv.Itoa(transaction.ID), transaction.PaymentType, int64(transaction.Amount))
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		transaction.VirtualAccount = req
		update := svc.model.Update(transaction)
		if update == -1 {
			return nil, err
		}
		resTransaction.PaymentType = "Bank Transfer"
		resTransaction.VirtualAccount = req
	case "10":
		req, err := svc.mtRequest.CreateTransactionGopay(strconv.Itoa(transaction.ID), transaction.PaymentType, int64(transaction.Amount))
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}
		transaction.UrlCallback = req
		update := svc.model.Update(transaction)
		if update == -1 {
			return nil, err
		}
		resTransaction.PaymentType = "Gopay"
		resTransaction.UrlCallback = req
	default:
		req, err := svc.mtRequest.CreateTransactionBank(strconv.Itoa(transaction.ID), transaction.PaymentType, int64(transaction.Amount))
		if err != nil {
			return nil, err
		}
		update := svc.model.Update(transaction)
		if update == -1 {
			return nil, err
		}
		resTransaction.PaymentType = "Bank Transfer"
		resTransaction.VirtualAccount = req
	}

	return &resTransaction, nil
}

func (svc *service) Modify(transactionData dtos.InputTransaction, transactionID int) bool {
	newTransaction := transaction.Transaction{}

	err := smapping.FillStruct(&newTransaction, smapping.MapFields(transactionData))
	if err != nil {
		log.Error(err)
		return false
	}

	newTransaction.ID = transactionID
	rowsAffected := svc.model.Update(newTransaction)

	if rowsAffected <= 0 {
		log.Error("There is No Transaction Updated!")
		return false
	}

	return true
}

func (svc *service) Remove(transactionID int) bool {
	rowsAffected := svc.model.DeleteByID(transactionID)

	if rowsAffected <= 0 {
		log.Error("There is No Transaction Deleted!")
		return false
	}

	return true
}

func (svc *service) Notifications(notificationPayload map[string]any) error {
	transactionID, exist := notificationPayload["order_id"].(string)
	if !exist {
		return errors.New("invalid notification payload")
	}

	transactionStatusResp, err := svc.coreAPIClient.CheckTransaction(transactionID)
	if err != nil {
		return err
	} else {
		if transactionStatusResp != nil {
			var status = svc.mtRequest.TransactionStatus(transactionStatusResp)
			transactionIDInt, err := strconv.Atoi(transactionID)
			if err != nil {
				return err
			}
			transaction := svc.model.SelectByID(transactionIDInt)
			transaction.Status = status.Order

			if transaction.Status == "5" {
				transaction.PaidAt = time.Now().Format("2006-01-02 15:04:05")
				update := svc.model.Update(*transaction)
				if update == -1 {
					return nil
				}
			} else {
				update := svc.model.Update(*transaction)
				if update == -1 {
					return nil
				}
			}

		}
	}

	return nil
}
