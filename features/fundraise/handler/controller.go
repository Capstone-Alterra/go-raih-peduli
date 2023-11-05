package handler

import (
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
	return &controller{
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetFundraises() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)

		page := pagination.Page
		size := pagination.Size
		title := ctx.QueryParam("title")

		if page <= 0 || size <= 0 {
			return ctx.JSON(400, helper.Response("Please provide query `page` and `size` in number!"))
		}

		fundraises := ctl.service.FindAll(page, size, title)

		if fundraises == nil {
			return ctx.JSON(404, helper.Response("There is No Fundraises!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": fundraises,
		}))
	}
}

func (ctl *controller) FundraiseDetails() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		fundraiseID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		fundraise := ctl.service.FindByID(fundraiseID)
		if fundraise == nil {
			return ctx.JSON(404, helper.Response("Fundraise Not Found!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": fundraise,
		}))
	}
}

func (ctl *controller) DeleteFundraise() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		fundraiseID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		fundraise := ctl.service.FindByID(fundraiseID)

		if fundraise == nil {
			return ctx.JSON(404, helper.Response("Fundraise Not Found!"))
		}

		delete := ctl.service.Remove(fundraiseID)

		if !delete {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Fundraise Success Deleted!", nil))
	}
}
