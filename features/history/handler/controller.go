package handler

import (
	"raihpeduli/features/history"
	helper "raihpeduli/helpers"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service history.Usecase
}

func New(service history.Usecase) history.Handler {
	return &controller{
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetHistoryFundraiseCreatedByUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var userId int

		if ctx.Get("user_id") != nil {
			userId = ctx.Get("user_id").(int)
		}

		fundraises, err := ctl.service.FindAllHistoryFundraiseCreatedByUser(userId)

		if err != nil {
			return ctx.JSON(500, helper.Response("something happend error"))
		}
		if len(fundraises) == 0 {
			return ctx.JSON(404, helper.Response("history fundraises not found"))
		}
		return ctx.JSON(200, helper.Response(
			"success", map[string]any{
				"data": fundraises,
			}))
	}
}

func (ctl *controller) GetHistoryVolunteerVacanciesCreatedByUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return nil
	}
}

func (ctrl *controller) GetHistoryVolunteerVacanciewsRegisterByUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := 0
		if ctx.Get("user_id") != nil {
			userID = ctx.Get("user_id").(int)
		}

		vacancies, err := ctrl.service.FindAllHistoryVolunteerVacanciewsRegisterByUser(userID)
		if err != nil {
			return ctx.JSON(500, helper.Response("something happend error"))
		}

		if vacancies == nil {
			return ctx.JSON(404, helper.Response("history volunteer vacancies not found"))

		}
		return ctx.JSON(200, helper.Response(
			"success", map[string]any{
				"data": vacancies,
			}))
	}
}

func (ctl *controller) GetHistoryUserTransaction() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		return nil
	}
}
