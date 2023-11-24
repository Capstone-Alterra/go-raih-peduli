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

func (ctl *controller) GetVacancies() echo.HandlerFunc {
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

		volunteers, totalData := ctl.service.FindAllVacancies(page, size, searchAndFilter)

		if volunteers == nil {
			return ctx.JSON(404, helpers.Response("There is No Volunteers!"))
		}

		paginationResponse := helpers.PaginationResponse(page, size, int(totalData))

		return ctx.JSON(200, helpers.Response("Success!", map[string]any{
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
			return ctx.JSON(404, helpers.Response("Volunteer Not Found!"))
		}

		return ctx.JSON(200, helpers.Response("Success!", map[string]any{
			"data": volunteer,
		}))
	}
}

func (ctl *controller) UpdateVacancy() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputVacancy{}

		volunteerID, errParam := strconv.Atoi(ctx.Param("id"))

		if errParam != nil {
			return ctx.JSON(400, helpers.Response(errParam.Error()))
		}

		volunteer := ctl.service.FindVacancyByID(volunteerID)

		if volunteer == nil {
			return ctx.JSON(404, helpers.Response("Volunteer Not Found!"))
		}

		ctx.Bind(&input)

		// validate = validator.New(validator.WithRequiredStructEnabled())
		// err := validate.Struct(input)

		// if err != nil {
		// 	errMap := helperss.ErrorMapValidation(err)
		// 	return ctx.JSON(400, helpers.Response("Bad Request!", map[string]any{
		// 		"error": errMap,
		// 	}))
		// }

		update := ctl.service.ModifyVacancy(input, volunteerID)

		if !update {
			return ctx.JSON(500, helpers.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helpers.Response("Volunteer Success Updated!"))
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
			return ctx.JSON(404, helpers.Response("Volunteer Not Found!"))
		}

		delete := ctl.service.RemoveVacancy(volunteerID)

		if !delete {
			return ctx.JSON(500, helpers.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helpers.Response("Volunteer Success Deleted!", nil))
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
			return ctx.JSON(500, helpers.Response("Controller : Something when wrong!", nil))
		}

		if err != nil {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("Succes", map[string]any{
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
			return ctx.JSON(500, helpers.Response("Controller : Something when wrong!"))
		}

		return ctx.JSON(200, helpers.Response("Apply Volunteer Success!", nil))
	}
}

func (ctl *controller) UpdateStatusRegistrar() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.StatusRegistrar{}

		ctx.Bind(&input)

		registrarID, errParam := strconv.Atoi(ctx.Param("id"))

		if errParam != nil {
			return ctx.JSON(400, helpers.Response(errParam.Error()))
		}

		update := ctl.service.UpdateStatusRegistrar(input.Status, registrarID)

		if !update {
			return ctx.JSON(500, helpers.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helpers.Response("Registrar Success Updated!"))
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
			return ctx.JSON(404, helpers.Response("There is No Volunteers!"))
		}

		paginationResponse := helpers.PaginationResponse(page, size, int(totalData))

		return ctx.JSON(200, helpers.Response("Success!", map[string]any{
			"data":       volunteers,
			"pagination": paginationResponse,
		}))
	}
}
