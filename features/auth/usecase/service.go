package usecase

import (
	"errors"
	"raihpeduli/features/auth"
	"raihpeduli/features/auth/dtos"
	"raihpeduli/helpers"
	"strconv"
)

type service struct {
	model auth.Repository
	jwt   helpers.JWTInterface
	hash  helpers.HashInterface
}

func New(model auth.Repository, jwt helpers.JWTInterface, hash helpers.HashInterface) auth.Usecase {
	return &service{
		model: model,
		jwt:   jwt,
		hash:  hash,
	}
}

func (svc *service) Login(data dtos.RequestLogin) (*dtos.LoginResponse, error) {
	user, err := svc.model.Login(data.Email)
	if err != nil {
		return nil, err
	}

	if !svc.hash.CompareHash(data.Password, user.Password) {
		return nil, errors.New("invalid password")
	}

	var fullname string
	if user.RoleID == 1 {
		fullname, err = svc.model.GetNameCustomer(user.ID)
	} else {
		fullname, err = svc.model.GetNameAdmin(user.ID)
	}
	if err != nil {
		return nil, err
	}

	tokenData := svc.jwt.GenerateJWT(strconv.Itoa(user.ID), strconv.Itoa(user.RoleID))
	return &dtos.LoginResponse{
		Name:         fullname,
		Email:        user.Email,
		Role:         user.RoleID,
		AccessToken:  tokenData["access_token"].(string),
		RefreshToken: tokenData["refresh_token"].(string),
	}, nil
}
