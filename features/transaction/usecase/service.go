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
	"github.com/sirupsen/logrus"
)

type service struct {
	model         transaction.Repository
	generator     helpers.GeneratorInterface
	mtRequest     helpers.MidtransInterface
	coreAPIClient coreapi.Client
	validation    helpers.ValidationInterface
}

func New(model transaction.Repository, generator helpers.GeneratorInterface, mtRequest helpers.MidtransInterface, coreAPIClient coreapi.Client, validation helpers.ValidationInterface) transaction.Usecase {
	return &service{
		model:         model,
		generator:     generator,
		mtRequest:     mtRequest,
		coreAPIClient: coreAPIClient,
		validation:    validation,
	}
}

func (svc *service) FindAll(page, size, roleID, userID int, keyword string) ([]dtos.ResTransaction, int64) {
	var transactions []dtos.ResTransaction
	var transactionData []transaction.Transaction
	var totalData int64
	if roleID == 2 {
		transactionData = svc.model.Paginate(page, size, keyword)
		totalData = svc.model.GetTotalData(keyword)
	} else {
		transactionData = svc.model.PaginateUser(page, size, userID)
		totalData = svc.model.GetTotalDataByUser(userID, keyword)
	}

	for _, transaction := range transactionData {
		var data dtos.ResTransaction

		if err := smapping.FillStruct(&data, smapping.MapFields(transaction)); err != nil {
			log.Error(err.Error())
		}

		data.Fullname = transaction.User.Fullname
		data.Address = transaction.User.Address
		data.PhoneNumber = transaction.User.PhoneNumber
		data.ProfilePicture = transaction.User.ProfilePicture
		data.Email = transaction.User.Email

		switch transaction.Status {
		case "2":
			data.Status = "Waiting For Payment"
		case "3":
			data.Status = "Failed / Cancelled"
		case "4":
			data.Status = "Transaction Success"
		case "5":
			data.Status = "Paid"
		default:
			data.Status = "Created"
		}

		switch transaction.PaymentType {
		case "4":
			data.PaymentType = "Bank Permata"
		case "5":
			data.PaymentType = "Bank CIMB"
		case "6":
			data.PaymentType = "Bank BCA"
		case "7":
			data.PaymentType = "Bank BRI"
		case "8":
			data.PaymentType = "Bank BNI"
		case "10":
			data.PaymentType = "Gopay"
		case "11":
			data.PaymentType = "Qris"
		default:
			data.PaymentType = "Other"
		}

		transactions = append(transactions, data)
	}

	return transactions, totalData
}

func (svc *service) FindByID(transactionID int) *dtos.ResTransaction {
	res := dtos.ResTransaction{}
	transaction := svc.model.SelectByID(transactionID)

	if transaction == nil {
		return nil
	}
	switch transaction.Status {
	case "2":
		transaction.Status = "Waiting For Payment"
	case "3":
		transaction.Status = "Failed / Cancelled"
	case "4":
		transaction.Status = "Transaction Success"
	case "5":
		transaction.Status = "Paid"
	default:
		transaction.Status = "Created"
	}

	switch transaction.PaymentType {
	case "4":
		transaction.PaymentType = "Bank Permata"
	case "5":
		transaction.PaymentType = "Bank CIMB"
	case "6":
		transaction.PaymentType = "Bank BCA"
	case "7":
		transaction.PaymentType = "Bank BRI"
	case "8":
		transaction.PaymentType = "Bank BNI"
	case "10":
		transaction.PaymentType = "Gopay"
	case "11":
		transaction.PaymentType = "Qris"
	default:
		transaction.PaymentType = "Other"
	}

	res.Fullname = transaction.User.Fullname
	res.Address = transaction.User.Address
	res.PhoneNumber = transaction.User.PhoneNumber
	res.ProfilePicture = transaction.User.ProfilePicture
	res.Email = transaction.User.Email

	err := smapping.FillStruct(&res, smapping.MapFields(transaction))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}
func (svc *service) Create(userID int, newTransaction dtos.InputTransaction) (*dtos.ResTransaction, error, []string) {
	if errMap := svc.validation.ValidateRequest(newTransaction); errMap != nil {
		return nil, nil, errMap
	}

	transaction := transaction.Transaction{}
	resTransaction := dtos.ResTransaction{}

	if newTransaction.Amount < 10000 {
		return nil, errors.New("Minimum donation ammount is Rp. 10.000"), nil
	}

	err := smapping.FillStruct(&transaction, smapping.MapFields(newTransaction))
	if err != nil {
		log.Error(err)
		return nil, err, nil
	}

	checkFundraise, err := svc.model.CountByID(newTransaction.FundraiseID)
	if err != nil {
		return nil, err, nil
	}

	if checkFundraise <= 0 {
		return nil, errors.New("Fundraise not found"), nil
	}

	user := svc.model.SelectUserByID(userID)

	transaction.UserID = userID
	transaction.ID = svc.generator.GenerateRandomID()
	transaction.Status = "2"

	switch transaction.PaymentType {
	case "4", "5", "6", "7", "8", "9":
		req, validUntil, err := svc.mtRequest.CreateTransactionBank(strconv.Itoa(transaction.ID), transaction.PaymentType, int64(transaction.Amount))

		if err != nil {
			log.Error(err.Error())
			return nil, err, nil
		}
		transaction.ValidUntil = validUntil

		transactionID := svc.model.Insert(transaction)
		if transactionID == -1 {
			return nil, err, nil
		}
		transaction.VirtualAccount = req
		update := svc.model.Update(transaction)
		if update == -1 {
			return nil, err, nil
		}
		resTransaction.Fullname = transaction.User.Fullname
		resTransaction.PaymentType = "Bank Transfer"
		resTransaction.VirtualAccount = req
		resTransaction.ID = transaction.ID
		resTransaction.Amount = int(transaction.Amount)
		resTransaction.Status = "Created"
		resTransaction.ValidUntil = validUntil
		resTransaction.UserID = userID
		resTransaction.Fullname = user.Fullname
		resTransaction.Address = user.Address
		resTransaction.PhoneNumber = user.PhoneNumber
		resTransaction.ProfilePicture = user.ProfilePicture
		resTransaction.FundraiseID = transaction.FundraiseID
	case "10":
		req, validUntil, err := svc.mtRequest.CreateTransactionGopay(strconv.Itoa(transaction.ID), transaction.PaymentType, int64(transaction.Amount))
		if err != nil {
			log.Error(err.Error())
			return nil, err, nil
		}
		transaction.ValidUntil = validUntil
		transactionID := svc.model.Insert(transaction)
		if transactionID == -1 {
			return nil, err, nil
		}
		transaction.UrlCallback = req
		update := svc.model.Update(transaction)
		if update == -1 {
			return nil, err, nil
		}
		resTransaction.PaymentType = "Gopay"
		resTransaction.UrlCallback = req
		resTransaction.ID = transaction.ID
		resTransaction.Amount = int(transaction.Amount)
		resTransaction.Status = "Created"
		resTransaction.ValidUntil = validUntil
		resTransaction.UserID = userID
		resTransaction.Fullname = user.Fullname
		resTransaction.Address = user.Address
		resTransaction.PhoneNumber = user.PhoneNumber
		resTransaction.ProfilePicture = user.ProfilePicture
		resTransaction.FundraiseID = transaction.FundraiseID
	case "11":
		req, validUntil, err := svc.mtRequest.CreateTransactionQris(strconv.Itoa(transaction.ID), transaction.PaymentType, int64(transaction.Amount))
		if err != nil {
			log.Error(err.Error())
			return nil, err, nil
		}
		transaction.ValidUntil = validUntil
		transactionID := svc.model.Insert(transaction)
		if transactionID == -1 {
			return nil, err, nil
		}
		transaction.UrlCallback = req
		update := svc.model.Update(transaction)
		if update == -1 {
			return nil, err, nil
		}
		resTransaction.PaymentType = "Qris"
		resTransaction.UrlCallback = req
		resTransaction.ID = transaction.ID
		resTransaction.Amount = int(transaction.Amount)
		resTransaction.Status = "Created"
		resTransaction.ValidUntil = validUntil
		resTransaction.UserID = userID
		resTransaction.Fullname = user.Fullname
		resTransaction.Address = user.Address
		resTransaction.PhoneNumber = user.PhoneNumber
		resTransaction.ProfilePicture = user.ProfilePicture
		resTransaction.FundraiseID = transaction.FundraiseID
	default:
		req, validUntil, err := svc.mtRequest.CreateTransactionBank(strconv.Itoa(transaction.ID), transaction.PaymentType, int64(transaction.Amount))
		if err != nil {
			return nil, err, nil
		}
		transaction.ValidUntil = validUntil
		transactionID := svc.model.Insert(transaction)
		if transactionID == -1 {
			return nil, err, nil
		}
		update := svc.model.Update(transaction)
		if update == -1 {
			return nil, err, nil
		}
		resTransaction.PaymentType = "Bank Transfer"
		resTransaction.VirtualAccount = req
		resTransaction.ID = transaction.ID
		resTransaction.Amount = int(transaction.Amount)
		resTransaction.Status = "Created"
		resTransaction.ValidUntil = validUntil
		resTransaction.UserID = userID
		resTransaction.Fullname = user.Fullname
		resTransaction.Address = user.Address
		resTransaction.PhoneNumber = user.PhoneNumber
		resTransaction.ProfilePicture = user.ProfilePicture
		resTransaction.FundraiseID = transaction.FundraiseID
	}

	return &resTransaction, nil, nil
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
				go func() {
					switch transaction.PaymentType {
					case "4":
						logrus.Println(transaction.User.Email)
						err := svc.model.SendPaymentConfirmation(transaction.User.Email, transaction.Amount, transaction.FundraiseID, "Bank Permata")
						if err != nil {
							logrus.Println(err.Error())
						}
					case "5":
						logrus.Println(transaction.User.Email)
						err := svc.model.SendPaymentConfirmation(transaction.User.Email, transaction.Amount, transaction.FundraiseID, "Bank CIMB")
						if err != nil {
							logrus.Println(err.Error())
						}
					case "6":
						logrus.Println(transaction.User.Email)
						err := svc.model.SendPaymentConfirmation(transaction.User.Email, transaction.Amount, transaction.FundraiseID, "Bank BCA")
						if err != nil {
							logrus.Println(err.Error())
						}
					case "7":
						logrus.Println(transaction.User.Email)
						err := svc.model.SendPaymentConfirmation(transaction.User.Email, transaction.Amount, transaction.FundraiseID, "Bank BRI")
						if err != nil {
							logrus.Println(err.Error())
						}
					case "8":
						logrus.Println(transaction.User.Email)
						err := svc.model.SendPaymentConfirmation(transaction.User.Email, transaction.Amount, transaction.FundraiseID, "Bank BNI")
						if err != nil {
							logrus.Println(err.Error())
						}
					case "10":
						logrus.Println(transaction.User.Email)
						err := svc.model.SendPaymentConfirmation(transaction.User.Email, transaction.Amount, transaction.FundraiseID, "Gopay")
						if err != nil {
							logrus.Println(err.Error())
						}
					case "11":
						logrus.Println(transaction.User.Email)
						err := svc.model.SendPaymentConfirmation(transaction.User.Email, transaction.Amount, transaction.FundraiseID, "Qris")
						if err != nil {
							logrus.Println(err.Error())
						}
					default:
						transaction.PaymentType = "Other"
					}

				}()

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
