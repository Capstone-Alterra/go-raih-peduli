package usecase

import (
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
	hash  helpers.HashInterface
}

func New(model admin.Repository, jwt helpers.JWTInterface, hash helpers.HashInterface) admin.Usecase {
	return &service{
		model: model,
		jwt:   jwt,
		hash:  hash,
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

func (svc *service) Create(newData dtos.InputAdmin) (*dtos.ResAdmin, error) {
	newUser := admin.User{}
	newAdmin := admin.Admin{}

	err := smapping.FillStruct(&newUser, smapping.MapFields(newData))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = smapping.FillStruct(&newAdmin, smapping.MapFields(newData))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	newUser.Password = svc.hash.HashPassword(newUser.Password)
	user, err := svc.model.InsertUser(&newUser)
	if err != nil {
		return nil, err
	}

	newAdmin.UserID = user.ID
	admin, err := svc.model.InsertAdmin(&newAdmin)
	if err != nil {
		return nil, err
	}

	resAdmin := dtos.ResAdmin{}
	resAdmin.RoleID = user.RoleID
	resAdmin.Email = user.Email
	errRes := smapping.FillStruct(&resAdmin, smapping.MapFields(admin))
	if errRes != nil {
		log.Error(errRes)
		return nil, err
	}

	userID := strconv.Itoa(resAdmin.UserID)
	roleID := strconv.Itoa(resAdmin.RoleID)
	tokenData := svc.jwt.GenerateJWT(userID, roleID)

	if tokenData == nil {
		log.Error("Token process failed")
	}

	resAdmin.AccessToken = tokenData["access_token"].(string)
	resAdmin.RefreshToken = tokenData["refresh_token"].(string)

	return &resAdmin, nil
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
