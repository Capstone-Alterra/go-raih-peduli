package usecase

import (
	"raihpeduli/features/volunteer"
	"raihpeduli/features/volunteer/dtos"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
)

type service struct {
	model volunteer.Repository
}

func New(model volunteer.Repository) volunteer.Usecase {
	return &service {
		model: model,
	}
}

func (svc *service) FindAll(page, size int) []dtos.ResVolunteer {
	var volunteers []dtos.ResVolunteer

	volunteersEnt := svc.model.Paginate(page, size)

	for _, volunteer := range volunteersEnt {
		var data dtos.ResVolunteer

		if err := smapping.FillStruct(&data, smapping.MapFields(volunteer)); err != nil {
			log.Error(err.Error())
		} 
		
		volunteers = append(volunteers, data)
	}

	return volunteers
}

func (svc *service) FindByID(volunteerID int) *dtos.ResVolunteer {
	res := dtos.ResVolunteer{}
	volunteer := svc.model.SelectByID(volunteerID)

	if volunteer == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(volunteer))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Modify(volunteerData dtos.InputVolunteer, volunteerID int) bool {
	newVolunteer := volunteer.VolunteerVacancies{}

	err := smapping.FillStruct(&newVolunteer, smapping.MapFields(volunteerData))
	if err != nil {
		log.Error(err)
		return false
	}

	newVolunteer.ID = volunteerID
	rowsAffected := svc.model.Update(newVolunteer)

	if rowsAffected <= 0 {
		log.Error("There is No Volunteer Updated!")
		return false
	}
	
	return true
}

func (svc *service) Remove(volunteerID int) bool {
	rowsAffected := svc.model.DeleteByID(volunteerID)

	if rowsAffected <= 0 {
		log.Error("There is No Volunteer Deleted!")
		return false
	}

	return true
}