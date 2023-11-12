package handler

import (
	"raihpeduli/helpers"
	helper "raihpeduli/helpers"
	"strconv"

	"raihpeduli/features/fundraise"
	"raihpeduli/features/fundraise/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service fundraise.Usecase
}

func New(service fundraise.Usecase) fundraise.Handler {
	return &controller {
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetFundraises() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)
		
		page := pagination.Page
		size := pagination.Size
		title := ctx.QueryParam("title")

		if page <= 0 || size <= 0 {
			page = 1
			size = 10
		}

		fundraises := ctl.service.FindAll(page, size, title)

		if fundraises == nil {
			return ctx.JSON(404, helper.Response("fundraises not found"))
		}

		return ctx.JSON(200, helper.Response("success", map[string]any {
			"data": fundraises,
		}))
	}
}


func (ctl *controller) FundraiseDetails() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		fundraiseID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		fundraise := ctl.service.FindByID(fundraiseID)

		if fundraise == nil {
			return ctx.JSON(404, helper.Response("fundraise not found"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": fundraise,
		}))
	}
}

func (ctl *controller) CreateFundraise() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		input := dtos.InputFundraise{}

		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())

		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any {
				"error": errMap,
			}))
		}

		userID := ctx.Get("user_id")

		fundraise, err := ctl.service.Create(input, userID.(int))

		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": fundraise,
		}))
	}
}

func (ctl *controller) UpdateFundraise() echo.HandlerFunc {
	return func (ctx echo.Context) error {
		input := dtos.InputFundraise{}

		fundraiseID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		fundraise := ctl.service.FindByID(fundraiseID)

		if fundraise == nil {
			return ctx.JSON(404, helper.Response("fundraise not found"))
		}
		
		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())
		

		if err := validate.Struct(input); err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any {
				"error": errMap,
			}))
		}

		update := ctl.service.Modify(input, fundraiseID)

		if !update {
			return ctx.JSON(500, helper.Response("something went wrong"))
		}

		return ctx.JSON(200, helper.Response("Fundraise Success Updated!"))
	}
}

func (ctl *controller) DeleteFundraise() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		fundraiseID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		fundraise := ctl.service.FindByID(fundraiseID)

		if fundraise == nil {
			return ctx.JSON(404, helper.Response("fundraise not found"))
		}

		delete := ctl.service.Remove(fundraiseID)

		if !delete {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Fundraise Success Deleted!", nil))
	}
}
