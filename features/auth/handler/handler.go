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

func (ctl *controller) Login() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		loginData := dtos.RequestLogin{}

		if err := ctx.Bind(&loginData); err != nil {
			return ctx.JSON(400, helpers.Response("invalid request body"))
		}

		loginRes, err := ctl.service.Login(loginData)
		if err != nil {
			return ctx.JSON(401, helpers.Response("invalid credentials"))
		}

		return ctx.JSON(200, helpers.Response("success", map[string]any{
			"data": loginRes,
		}))
	}
}

func (ctl *controller) RegisterUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputUser{}

		ctx.Bind(&input)

		user, errMap, err := ctl.service.Register(input)
		if errMap != nil {
			return ctx.JSON(400, helpers.Response("missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(400, helpers.Response("bad request", map[string]any{
				"error": err.Error(),
			}))
		}

		return ctx.JSON(200, helpers.Response("success", map[string]any{
			"data": user,
		}))
	}
}

func (ctl *controller) ResendOTP() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.ResendOTP{}

		ctx.Bind(&input)

		result := ctl.service.ResendOTP(input.Email)
		if !result {
			return ctx.JSON(500, helpers.Response("something went wrong"))
		}

		return ctx.JSON(200, helpers.Response("OTP has been sent via email"))
	}
}

func (ctl *controller) RefreshJWT() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		jwt := dtos.RefreshJWT{}
		ctx.Bind(&jwt)

		refershJWT, err := ctl.service.RefreshJWT(jwt)
		if err != nil {
			if err.Error() == "validate token failed" {
				return ctx.JSON(400, helpers.Response("invalid jwt token"))
			}

			return ctx.JSON(500, helpers.Response("something went wrong"))
		}

		return ctx.JSON(200, helpers.Response("success", map[string]any{
			"data": refershJWT,
		}))
	}
}
