package handler

import (
	"mime/multipart"
	"raihpeduli/helpers"
	"strconv"

	"raihpeduli/features/volunteer"
	"raihpeduli/features/volunteer/dtos"

	"github.com/labstack/echo/v4"
)

type controller struct {
	service volunteer.Usecase
}

func New(service volunteer.Usecase) volunteer.Handler {
	return &controller{
		service: service,
	}
}

func (ctl *controller) GetVacancies(suffix string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)

		if pagination.Page < 1 || pagination.PageSize < 1 {
			pagination.Page = 1
			pagination.PageSize = 20
		}

		searchAndFilter := dtos.SearchAndFilter{}
		ctx.Bind(&searchAndFilter)

		page := pagination.Page
		size := pagination.PageSize

		volunteers, totalData := ctl.service.FindAllVacancies(page, size, searchAndFilter, suffix)

		if volunteers == nil {
			return ctx.JSON(404, helpers.Response("there is no volunteers!"))
		}

		paginationResponse := helpers.PaginationResponse(page, size, int(totalData))

		return ctx.JSON(200, helpers.Response("success!", map[string]any{
			"data":       volunteers,
			"pagination": paginationResponse,
		}))
	}
}

func (ctl *controller) VacancyDetails() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		volunteerID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helpers.Response(err.Error()))
		}

		volunteer := ctl.service.FindVacancyByID(volunteerID)

		if volunteer == nil {
			return ctx.JSON(404, helpers.Response("volunteer not found!"))
		}

		return ctx.JSON(200, helpers.Response("success!", map[string]any{
			"data": volunteer,
		}))
	}
}

func (ctl *controller) UpdateVacancy() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputVacancy{}

		vacancyID, errParam := strconv.Atoi(ctx.Param("id"))

		if errParam != nil {
			return ctx.JSON(400, helpers.Response(errParam.Error()))
		}

		vacancy := ctl.service.FindVacancyByID(vacancyID)

		if vacancy == nil {
			return ctx.JSON(404, helpers.Response("volunteer not found!"))
		}

		ctx.Bind(&input)

		fileHeader, err := ctx.FormFile("photo")
		var file multipart.File

		if err == nil {
			formFile, err := fileHeader.Open()

			if err != nil {
				return ctx.JSON(500, helpers.Response("something went wrong"))
			}

			file = formFile
		}

		result, errMap := ctl.service.ModifyVacancy(input, file, *vacancy)
		if errMap != nil {
			return ctx.JSON(400, helpers.Response("error missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if !result {
			return ctx.JSON(500, helpers.Response("something went wrong!"))
		}

		return ctx.JSON(200, helpers.Response("volunteer success updated!"))
	}
}

func (ctl *controller) DeleteVacancy() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		volunteerID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helpers.Response(err.Error()))
		}

		volunteer := ctl.service.FindVacancyByID(volunteerID)

		if volunteer == nil {
			return ctx.JSON(404, helpers.Response("volunteer not found!"))
		}

		delete := ctl.service.RemoveVacancy(volunteerID)

		if !delete {
			return ctx.JSON(500, helpers.Response("something went wrong!"))
		}

		return ctx.JSON(200, helpers.Response("volunteer success deleted!", nil))
	}
}

func (ctl *controller) CreateVacancy() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputVacancy{}

		ctx.Bind(&input)

		userID := ctx.Get("user_id")

		fileHeader, err := ctx.FormFile("photo")
		var file multipart.File

		if err == nil {
			formFile, err := fileHeader.Open()

			if err != nil {
				return ctx.JSON(500, helpers.Response("something went wrong"))
			}

			file = formFile
		}

		volun, errMap, err := ctl.service.CreateVacancy(input, userID.(int), file)

		if errMap != nil {
			return ctx.JSON(400, helpers.Response("missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if volun == nil {
			return ctx.JSON(500, helpers.Response("something when wrong!", nil))
		}

		if err != nil {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("succes", map[string]any{
			"data": volun,
		}))
	}
}

func (ctl *controller) ApplyVacancy() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.ApplyVacancy{}

		ctx.Bind(&input)

		userID := ctx.Get("user_id")

		fileHeader, err := ctx.FormFile("resume")
		var file multipart.File

		if err == nil {
			formFile, err := fileHeader.Open()

			if err != nil {
				return ctx.JSON(500, helpers.Response("something went wrong"))
			}

			file = formFile
		}

		result, errMap := ctl.service.RegisterVacancy(input, userID.(int), file)
		if errMap != nil {
			return ctx.JSON(400, helpers.Response("missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if !result {
			return ctx.JSON(500, helpers.Response("something when wrong!"))
		}

		return ctx.JSON(200, helpers.Response("apply volunteer success!", nil))
	}
}

func (ctl *controller) UpdateStatusRegistrar() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.StatusRegistrar{}

		ctx.Bind(&input)

		volunteerID, err := strconv.Atoi(ctx.Param("volunteer_id"))

		if err != nil {
			return ctx.JSON(400, helpers.Response(err.Error()))
		}

		update := ctl.service.UpdateStatusRegistrar(input.Status, volunteerID)

		if !update {
			return ctx.JSON(500, helpers.Response("something went wrong!"))
		}

		return ctx.JSON(200, helpers.Response("registrar success updated!"))
	}
}

func (ctl *controller) GetVolunteersByVacancyID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}

		vacancyID, errParam := strconv.Atoi(ctx.Param("vacancy_id"))
		if errParam != nil {
			return ctx.JSON(400, helpers.Response(errParam.Error()))
		}

		ctx.Bind(&pagination)

		if pagination.Page < 1 || pagination.PageSize < 1 {
			pagination.Page = 1
			pagination.PageSize = 20
		}

		page := pagination.Page
		size := pagination.PageSize
		name := ctx.QueryParam("name")

		volunteers, totalData := ctl.service.FindAllVolunteersByVacancyID(page, size, vacancyID, name)

		if volunteers == nil {
			return ctx.JSON(404, helpers.Response("there is no volunteers!"))
		}

		paginationResponse := helpers.PaginationResponse(page, size, int(totalData))

		return ctx.JSON(200, helpers.Response("success!", map[string]any{
			"data":       volunteers,
			"pagination": paginationResponse,
		}))
	}
}


func (ctl *controller) GetVolunteer() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		vacancyID, err := strconv.Atoi(ctx.Param("vacancy_id"))
		volunteerID, err := strconv.Atoi(ctx.Param("volunteer_id"))

		if err != nil {
			return ctx.JSON(400, helpers.Response(err.Error()))
		}

		volunteer := ctl.service.FindDetailVolunteers(vacancyID, volunteerID)

		if volunteer.Fullname == "" {
			return ctx.JSON(404, helpers.Response("volunteer not found!"))
		}

		return ctx.JSON(200, helpers.Response("success!", map[string]any{
			"data": volunteer,
		}))
	}
}