package usecase

import (
	"raihpeduli/features/_blueprint"
	"raihpeduli/features/_blueprint/dtos"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
)

type service struct {
	model _blueprint.Repository
}

func New(model _blueprint.Repository) _blueprint.Usecase {
	return &service {
		model: model,
	}
}

func (svc *service) FindAll(page, size int) []dtos.ResPlaceholder {
	var placeholders []dtos.ResPlaceholder

	placeholdersEnt := svc.model.Paginate(page, size)

	for _, placeholder := range placeholdersEnt {
		var data dtos.ResPlaceholder

		if err := smapping.FillStruct(&data, smapping.MapFields(placeholder)); err != nil {
			log.Error(err.Error())
		} 
		
		placeholders = append(placeholders, data)
	}

	return placeholders
}

func (svc *service) FindByID(placeholderID int) *dtos.ResPlaceholder {
	res := dtos.ResPlaceholder{}
	placeholder := svc.model.SelectByID(placeholderID)

	if placeholder == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(placeholder))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newPlaceholder dtos.InputPlaceholder) *dtos.ResPlaceholder {
	placeholder := _blueprint.Placeholder{}
	
	err := smapping.FillStruct(&placeholder, smapping.MapFields(newPlaceholder))
	if err != nil {
		log.Error(err)
		return nil
	}

	placeholderID := svc.model.Insert(placeholder)

	if placeholderID == -1 {
		return nil
	}

	resPlaceholder := dtos.ResPlaceholder{}
	errRes := smapping.FillStruct(&resPlaceholder, smapping.MapFields(newPlaceholder))
	if errRes != nil {
		log.Error(errRes)
		return nil
	}

	return &resPlaceholder
}

func (svc *service) Modify(placeholderData dtos.InputPlaceholder, placeholderID int) bool {
	newPlaceholder := _blueprint.Placeholder{}

	err := smapping.FillStruct(&newPlaceholder, smapping.MapFields(placeholderData))
	if err != nil {
		log.Error(err)
		return false
	}

	newPlaceholder.ID = placeholderID
	rowsAffected := svc.model.Update(newPlaceholder)

	if rowsAffected <= 0 {
		log.Error("There is No Placeholder Updated!")
		return false
	}
	
	return true
}

func (svc *service) Remove(placeholderID int) bool {
	rowsAffected := svc.model.DeleteByID(placeholderID)

	if rowsAffected <= 0 {
		log.Error("There is No Placeholder Deleted!")
		return false
	}

	return true
}