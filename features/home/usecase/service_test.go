package usecase

import (
	"raihpeduli/features/home/dtos"
	"raihpeduli/features/home/mocks"
	"testing"
)

func TestFindAll(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var service = New(repository)

	var resFundraises = []dtos.ResFundraise{
	}

	var resVolunteers = dtos.ResVolunteer{

	}

	var resNews = dtos.ResNews{

	}

	var resHome = dtos.ResGetHome{

	}
}