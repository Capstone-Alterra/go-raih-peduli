package handler

import (
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
		ctx.Bind(&pagination)

		page := pagination.Page
		size := pagination.Size
		title := ctx.QueryParam("title")
		skill := ctx.QueryParam("skill")

		if page <= 0 || size <= 0 {
			return ctx.JSON(400, helper.Response("Please provide query `page` and `size` in number!"))
		}

		volunteers := ctl.service.FindAll(page, size, title, skill)

		if volunteers == nil {
			return ctx.JSON(404, helper.Response("There is No Volunteers!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": volunteers,
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

		validate = validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any{
				"error": errMap,
			}))
		}

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
