package handler

import (
	helper "raihpeduli/helpers"

	"raihpeduli/features/chatbot"
	"raihpeduli/features/chatbot/dtos"

	"github.com/labstack/echo/v4"
)

type controller struct {
	service chatbot.Usecase
}

func New(service chatbot.Usecase) chatbot.Handler {
	return &controller {
		service: service,
	}
}

func (ctl *controller) GetChatHistory() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)

		userID := ctx.Get("user_id").(int)

		chatHistories := ctl.service.FindAllChat(userID)

		if chatHistories == nil || len(chatHistories) == 0 {
			return ctx.JSON(404, helper.Response("there is no chat history for this user"))
		}

		return ctx.JSON(200, helper.Response("success", map[string]any {
			"data": chatHistories,
		}))
	}
}

func (ctl *controller) GetNewsContentGeneration() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		input := dtos.InputMessage{}

		ctx.Bind(&input)

		message, errMap, err := ctl.service.SetContentForNews(input)

		if errMap != nil {
			return ctx.JSON(400, helper.Response("error missing some data", map[string]any {
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		return ctx.JSON(200, helper.Response("success", map[string]any {
			"data": message,
		}))
	}
}

func (ctl *controller) SendQuestion() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		input := dtos.InputMessage{}

		ctx.Bind(&input)

		userID := 0

		if ctx.Get("user_id") != nil {
			userID = ctx.Get("user_id").(int)
		}

		message, errMap, err := ctl.service.SetReplyMessage(input, userID)

		if errMap != nil {
			return ctx.JSON(400, helper.Response("error missing some data", map[string]any {
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		return ctx.JSON(200, helper.Response("success", map[string]any {
			"data": message,
		}))
	}
}

func (ctl *controller) DeleteChatHistory() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		userID := ctx.Get("user_id").(int)

		if err := ctl.service.ClearHistory(userID); err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		return ctx.JSON(200, helper.Response("success cleared chat history", nil))
	}
}
