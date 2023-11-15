package handler

import (
	"mime/multipart"
	helper "raihpeduli/helpers"
	"strconv"

	"raihpeduli/features/news"
	"raihpeduli/features/news/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service news.Usecase
}

func New(service news.Usecase) news.Handler {
	return &controller{
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetNews() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)

		page := pagination.Page
		size := pagination.Size
		keyword := ctx.QueryParam("title")

		if page <= 0 || size <= 0 {
			return ctx.JSON(400, helper.Response("Please provide query `page` and `size` in number!"))
		}

		newss := ctl.service.FindAll(page, size, keyword)

		if newss == nil {
			return ctx.JSON(404, helper.Response("There is No Newss!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": newss,
		}))
	}
}

func (ctl *controller) NewsDetails() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		newsID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		news := ctl.service.FindByID(newsID)

		if news == nil {
			return ctx.JSON(404, helper.Response("News Not Found!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": news,
		}))
	}
}

func (ctl *controller) CreateNews() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputNews{}

		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())

		if err := validate.Struct(input); err != nil {
			errMap := helper.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any{
				"error": errMap,
			}))
		}

		userID := ctx.Get("user_id")

		fileHeader, err := ctx.FormFile("photo")
		var file multipart.File

		if err == nil {
			formFile, err := fileHeader.Open()

			if err != nil {
				return ctx.JSON(500, helper.Response("something went wrong"))
			}

			file = formFile
		}
		news, err := ctl.service.Create(input, userID.(int), file)

		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error(), nil))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": news,
		}))
	}
}

func (ctl *controller) UpdateNews() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputNews{}

		newsID, errParam := strconv.Atoi(ctx.Param("id"))

		if errParam != nil {
			return ctx.JSON(400, helper.Response(errParam.Error()))
		}

		news := ctl.service.FindByID(newsID)

		if news == nil {
			return ctx.JSON(404, helper.Response("News Not Found!"))
		}

		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(input)

		if err != nil {
			errMap := helper.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any{
				"error": errMap,
			}))
		}

		update := ctl.service.Modify(input, newsID)

		if !update {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("News Success Updated!"))
	}
}

func (ctl *controller) DeleteNews() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		newsID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		news := ctl.service.FindByID(newsID)

		if news == nil {
			return ctx.JSON(404, helper.Response("News Not Found!"))
		}

		delete := ctl.service.Remove(newsID)

		if !delete {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("News Success Deleted!", nil))
	}
}
