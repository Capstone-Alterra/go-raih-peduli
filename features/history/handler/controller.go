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
			return ctx.JSON(500, helper.Response(err.Error()))
		}
		if fundraises == nil || len(fundraises) == 0 {
			return ctx.JSON(404, helper.Response("history fundraises created by users not found"))
		}
		return ctx.JSON(200, helper.Response(
			"success", map[string]any{
				"data": fundraises,
			}))
	}
}

func (ctl *controller) GetHistoryVolunteerVacanciesCreatedByUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := 0
		if ctx.Get("user_id") != nil {
			userID = ctx.Get("user_id").(int)
		}

		vacancies, err := ctl.service.FindAllHistoryVolunteerVacanciesCreatedByUser(userID)
		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		if vacancies == nil || len(vacancies) == 0 {
			return ctx.JSON(404, helper.Response("history volunteer vacancies created by users not found"))

		}
		return ctx.JSON(200, helper.Response(
			"success", map[string]any{
				"data": vacancies,
			}))
	}
}

func (ctrl *controller) GetHistoryVolunteerVacanciesRegisterByUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := 0
		if ctx.Get("user_id") != nil {
			userID = ctx.Get("user_id").(int)
		}

		vacancies, err := ctrl.service.FindAllHistoryVolunteerVacanciesRegisterByUser(userID)
		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		if vacancies == nil || len(vacancies) == 0 {
			return ctx.JSON(404, helper.Response("history volunteer vacancies registered by users not found"))

		}
		return ctx.JSON(200, helper.Response(
			"success", map[string]any{
				"data": vacancies,
			}))
	}
}

func (ctrl *controller) GetHistoryUserTransaction() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := 0
		if ctx.Get("user_id") != nil {
			userID = ctx.Get("user_id").(int)
		}

		donations, err := ctrl.service.FindAllHistoryUserTransaction(userID)
		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		if donations == nil || len(donations) == 0 {
			return ctx.JSON(404, helper.Response("donation not found"))

		}
		return ctx.JSON(200, helper.Response(
			"success", map[string]any{
				"data": donations,
			}))
	}
}

func (ctl *controller) GetAllHistory() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := 0
		if ctx.Get("user_id") != nil {
			userID = ctx.Get("user_id").(int)
		}
		response_data := make(map[string]any)
		fundraises, err := ctl.service.FindAllHistoryFundraiseCreatedByUser(userID)

		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}
		if len(fundraises) != 0 {
			// return ctx.JSON(404, helper.Response("history fundraises created by users not found"))
			response_data["created_fundraises"] = fundraises
		}

		vacancies, err := ctl.service.FindAllHistoryVolunteerVacanciesCreatedByUser(userID)
		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		if vacancies != nil || len(vacancies) != 0 {
			// return ctx.JSON(404, helper.Response("history volunteer vacancies  created by users not found"))
			response_data["volunteer_vacancies"] = vacancies
		}
		vacanciesReg, err := ctl.service.FindAllHistoryVolunteerVacanciesRegisterByUser(userID)
		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		if vacanciesReg != nil || len(vacanciesReg) != 0 {
			// return ctx.JSON(404, helper.Response("history volunteer vacancies registered by users not found"))
			response_data["volunteer_vacancies_registered"] = vacanciesReg
		}

		donations, err := ctl.service.FindAllHistoryUserTransaction(userID)
		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		if donations != nil || len(donations) != 0 {
			// return ctx.JSON(404, helper.Response("donation not found"))
			response_data["donations"] = donations
		}

		if response_data == nil || len(response_data) == 0 {
			return ctx.JSON(404, helper.Response("history not found"))
		}
		data := make(map[string]any)
		data["data"] = response_data
		return ctx.JSON(200, helper.Response("success", data))
	}
}
