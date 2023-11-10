package hendler

import (
	"raihpeduli/features/auth"
	"raihpeduli/features/auth/dtos"
	"raihpeduli/helpers"

	"github.com/labstack/echo/v4"
)

type controller struct {
	service auth.Usecase
}

func New(service auth.Usecase) auth.Handler {
	return &controller{
		service: service,
	}
}

func (ctl *controller) LoginCustomer() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		loginData := dtos.RequestLogin{}

		if err := ctx.Bind(&loginData); err != nil {
			return ctx.JSON(400, helpers.Response("Invalid request body!"))
		}

		loginRes, err := ctl.service.Login(loginData)
		if err != nil {
			return ctx.JSON(401, helpers.Response("Invalid credentials!"))
		}

		return ctx.JSON(200, helpers.Response("Success!", map[string]any{
			"data": loginRes,
		}))
	}
}
