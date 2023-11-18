package handler

import (
	"mime/multipart"
	"raihpeduli/helpers"
	helper "raihpeduli/helpers"
	"strconv"

	"raihpeduli/features/volunteer"
	"raihpeduli/features/volunteer/dtos"

	"github.com/go-playground/validator/v10"
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

var validate *validator.Validate

func (ctl *controller) GetVolunteers() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}
		paginationResponse := dtos.PaginationResponse{}

		ctx.Bind(&pagination)

		if pagination.Page < 1 || pagination.Size < 1 {
			pagination.Page = 1
			pagination.Size = 20
		}

		page := pagination.Page
		size := pagination.Size
		title := ctx.QueryParam("title")
		skill := ctx.QueryParam("skill")
		city := ctx.QueryParam("city")

		volunteers, totalData := ctl.service.FindAll(page, size, title, skill, city)

		if volunteers == nil {
			return ctx.JSON(404, helper.Response("There is No Volunteers!"))
		}

		if pagination.Size >= int(totalData) {
			paginationResponse.PreviousPage = -1
			paginationResponse.NextPage = -1
		} else if pagination.Size < int(totalData) && pagination.Page == 1 {
			paginationResponse.PreviousPage = -1
			paginationResponse.NextPage = pagination.Page + 1
		} else {
			paginationResponse.PreviousPage = pagination.Page - 1
			paginationResponse.NextPage = pagination.Page + 1
		}

		paginationResponse.TotalData = totalData
		paginationResponse.CurrentPage = pagination.Page
		paginationResponse.TotalPage = (int(totalData) + pagination.Size - 1) / pagination.Size

		if paginationResponse.CurrentPage == paginationResponse.TotalPage {
			paginationResponse.NextPage = -1
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data":       volunteers,
			"pagination": paginationResponse,
		}))
	}
}

func (ctl *controller) VolunteerDetails() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		volunteerID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		volunteer := ctl.service.FindByID(volunteerID)

		if volunteer == nil {
			return ctx.JSON(404, helper.Response("Volunteer Not Found!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": volunteer,
		}))
	}
}

func (ctl *controller) UpdateVolunteer() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputVolunteer{}

		volunteerID, errParam := strconv.Atoi(ctx.Param("id"))

		if errParam != nil {
			return ctx.JSON(400, helper.Response(errParam.Error()))
		}

		volunteer := ctl.service.FindByID(volunteerID)

		if volunteer == nil {
			return ctx.JSON(404, helper.Response("Volunteer Not Found!"))
		}

		ctx.Bind(&input)

		// validate = validator.New(validator.WithRequiredStructEnabled())
		// err := validate.Struct(input)

		// if err != nil {
		// 	errMap := helpers.ErrorMapValidation(err)
		// 	return ctx.JSON(400, helper.Response("Bad Request!", map[string]any{
		// 		"error": errMap,
		// 	}))
		// }

		update := ctl.service.Modify(input, volunteerID)

		if !update {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Volunteer Success Updated!"))
	}
}

func (ctl *controller) DeleteVolunteer() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		volunteerID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		volunteer := ctl.service.FindByID(volunteerID)

		if volunteer == nil {
			return ctx.JSON(404, helper.Response("Volunteer Not Found!"))
		}

		delete := ctl.service.Remove(volunteerID)

		if !delete {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Volunteer Success Deleted!", nil))
	}
}

func (ctl *controller) CreateVolunteer() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputVolunteer{}

		ctx.Bind(&input)

		userID := ctx.Get("user_id")

		fileHeader, err := ctx.FormFile("photo")
		var file multipart.File

		if err == nil {
			formFile, err := fileHeader.Open()

			if err != nil {
				return ctx.JSON(500, helper.Response("something went wrong"))
			}

			file = formFile
		}

		volun, errMap, err := ctl.service.Create(input, userID.(int), file)

		if errMap != nil {
			return ctx.JSON(400, helper.Response("missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if volun == nil {
			return ctx.JSON(500, helpers.Response("Controller : Something when wrong!", nil))
		}

		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("Succes", map[string]any{
			"data": volun,
		}))
	}
}

func (ctl *controller) ApplyVacancies() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.ApplyVolunteer{}

		ctx.Bind(&input)

		userID := ctx.Get("user_id")

		fileHeader, err := ctx.FormFile("resume")
		var file multipart.File

		if err == nil {
			formFile, err := fileHeader.Open()

			if err != nil {
				return ctx.JSON(500, helper.Response("something went wrong"))
			}

			file = formFile
		}

		result, errMap := ctl.service.Register(input, userID.(int), file)
		if errMap != nil {
			return ctx.JSON(400, helper.Response("missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if !result {
			return ctx.JSON(500, helpers.Response("Controller : Something when wrong!"))
		}

		return ctx.JSON(200, helper.Response("Apply Volunteer Success!", nil))
	}
}
