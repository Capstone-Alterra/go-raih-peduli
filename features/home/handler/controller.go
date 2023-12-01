package handler

import (
	"raihpeduli/features/home"
	"raihpeduli/features/home/dtos"
	helper "raihpeduli/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service home.Usecase
}

func New(service home.Usecase) home.Handler {
	return &controller{
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetMobileLanding() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)

		page := 1
		size := 2

		homes := ctl.service.FindAll(page, size)

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": homes,
		}))
	}
}

func (ctl *controller) GetWebLanding() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)

		page := 1
		size := 5

		homes := ctl.service.FindAllWeb(page, size)

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": homes,
		}))
	}
}
