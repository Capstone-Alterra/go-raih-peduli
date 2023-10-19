package usecase

import (
	"errors"
	"raihpeduli/features/admin"
	"raihpeduli/features/admin/dtos"
	"raihpeduli/helpers"
	"strconv"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
)

type service struct {
	model admin.Repository
	jwt   helpers.JWTInterface
}

func New(model admin.Repository, jwt helpers.JWTInterface) admin.Usecase {
	return &service{
		model: model,
		jwt:   jwt,
	}
}

func (svc *service) FindAll(page, size int) []dtos.ResAdmin {
	var admins []dtos.ResAdmin

	adminsEnt := svc.model.Paginate(page, size)

	for _, admin := range adminsEnt {
		var data dtos.ResAdmin

		if err := smapping.FillStruct(&data, smapping.MapFields(admin)); err != nil {
			log.Error(err.Error())
		}

		admins = append(admins, data)
	}

	return admins
}

func (svc *service) FindByID(adminID int) *dtos.ResAdmin {
	res := dtos.ResAdmin{}
	admin := svc.model.SelectByID(adminID)

	if admin == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(admin))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newAdmin dtos.InputAdmin) *dtos.ResAdmin {
	admin := admin.Admin{}

	err := smapping.FillStruct(&admin, smapping.MapFields(newAdmin))
	if err != nil {
		log.Error(err)
		return nil
	}

	admin.Password = helpers.HashPassword(admin.Password)

	adminID := svc.model.Insert(&admin)

	if adminID == nil {
		return nil
	}

	resAdmin := dtos.ResAdmin{}
	errRes := smapping.FillStruct(&resAdmin, smapping.MapFields(admin))
	if errRes != nil {
		log.Error(errRes)
		return nil
	}

	ID := strconv.Itoa(resAdmin.ID)
	tokenData := svc.jwt.GenerateJWT(ID)

	if tokenData == nil {
		log.Error("Token process failed")
	}

	resAdmin.Token = tokenData

	return &resAdmin
}

func (svc *service) Modify(adminData dtos.InputAdmin, adminID int) bool {
	newAdmin := admin.Admin{}

	err := smapping.FillStruct(&newAdmin, smapping.MapFields(adminData))
	if err != nil {
		log.Error(err)
		return false
	}

	newAdmin.ID = adminID
	rowsAffected := svc.model.Update(newAdmin)

	if rowsAffected <= 0 {
		log.Error("There is No Admin Updated!")
		return false
	}

	return true
}

func (svc *service) Remove(adminID int) bool {
	rowsAffected := svc.model.DeleteByID(adminID)

	if rowsAffected <= 0 {
		log.Error("There is No Admin Deleted!")
		return false
	}

	return true
}

func (svc *service) Login(email, password string) (*dtos.ResLogin, error) {
	admin, err := svc.model.Login(email, password)
	if err != nil {
		return nil, err
	}

	if !helpers.CompareHash(password, admin.Password) {
		return nil, errors.New("invalid password")
	}

	tokenData := svc.jwt.GenerateJWT(strconv.Itoa(admin.ID))
	return &dtos.ResLogin{
		Name:  admin.Fullname,
		Email: admin.Email,
		Role:  "2",
		Token: tokenData,
	}, nil
}


