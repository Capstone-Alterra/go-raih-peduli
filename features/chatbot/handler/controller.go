package handler

import (
	"raihpeduli/helpers"
	helper "raihpeduli/helpers"
	"strconv"

	"raihpeduli/features/chatbot"
	"raihpeduli/features/chatbot/dtos"

	"github.com/go-playground/validator/v10"
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

var validate *validator.Validate

func (ctl *controller) GetChatHistory() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)
		
		page := pagination.Page
		size := pagination.Size

		if page <= 0 || size <= 0 {
			return ctx.JSON(400, helper.Response("Please provide query `page` and `size` in number!"))
		}

		chatbots := ctl.service.FindAll(page, size)

		if chatbots == nil {
			return ctx.JSON(404, helper.Response("There is No Chatbots!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": chatbots,
		}))
	}
}

func (ctl *controller) SendMessage() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		input := dtos.InputMessage{}

		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())

		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any {
				"error": errMap,
			}))
		}

		chatbot := ctl.service.Create(input)

		if chatbot == nil {
			return ctx.JSON(500, helper.Response("Something went Wrong!", nil))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": chatbot,
		}))
	}
}

func (ctl *controller) DeleteChatHistory() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		chatbotID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		chatbot := ctl.service.FindByID(chatbotID)

		if chatbot == nil {
			return ctx.JSON(404, helper.Response("Chatbot Not Found!"))
		}

		delete := ctl.service.Remove(chatbotID)

		if !delete {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Chatbot Success Deleted!", nil))
	}
}
