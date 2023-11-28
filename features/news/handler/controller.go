package handler

import (
	"mime/multipart"
	"raihpeduli/helpers"
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

		searchAndFilter := dtos.SearchAndFilter{}

		ctx.Bind(&searchAndFilter)

		userID := 0

		if ctx.Get("user_id") != nil {
			userID = ctx.Get("user_id").(int)
		}

		newss, totalData := ctl.service.FindAll(pagination, searchAndFilter, userID)

		if newss == nil {
			return ctx.JSON(404, helper.Response("news not found"))
		}

		page := pagination.Page
		pageSize := pagination.PageSize

		if page <= 0 || pageSize <= 0 {
			page = 1
			pageSize = 10
		}

		paginationResponse := helpers.PaginationResponse(page, pageSize, int(totalData))

		return ctx.JSON(200, helper.Response("success", map[string]any{
			"data":       newss,
			"pagination": paginationResponse,
		}))
	}
}

func (ctl *controller) NewsDetails() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		newsID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		userID := 0

		if ctx.Get("user_id") != nil {
			userID = ctx.Get("user_id").(int)
		}

		news := ctl.service.FindByID(newsID, userID)

		if news == nil {
			return ctx.JSON(404, helper.Response("news not found!"))
		}

		return ctx.JSON(200, helper.Response("success", map[string]any{
			"data": news,
		}))
	}
}

func (ctl *controller) CreateNews() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputNews{}

		ctx.Bind(&input)

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
		news, errMap, err := ctl.service.Create(input, userID.(int), file)

		if errMap != nil {
			return ctx.JSON(400, helper.Response("error missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error(), nil))
		}

		return ctx.JSON(200, helper.Response("success created fundraise", map[string]any{
			"data": news,
		}))
	}
}

func (ctl *controller) UpdateNews() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputNews{}

		newsID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		news := ctl.service.FindByID(newsID, 0)

		if news == nil {
			return ctx.JSON(404, helper.Response("news not found"))
		}

		ctx.Bind(&input)

		fileHeader, err := ctx.FormFile("photo")
		var file multipart.File

		if err == nil {
			formFile, err := fileHeader.Open()
			if err != nil {
				return ctx.JSON(500, helper.Response("something went wrong"))
			}

			file = formFile
		}

		errMap, err := ctl.service.Modify(input, file, *news)

		if errMap != nil {
			return ctx.JSON(500, helper.Response("something went wrong"))
		}

		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		return ctx.JSON(200, helper.Response("success updated news"))
	}

}

func (ctl *controller) DeleteNews() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		newsID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		news := ctl.service.FindByID(newsID, 0)

		if news == nil {
			return ctx.JSON(404, helper.Response("news not found"))
		}

		delete := ctl.service.Remove(newsID)

		if !delete {
			return ctx.JSON(500, helper.Response("something went wrong"))
		}

		return ctx.JSON(200, helper.Response("success deleted news", nil))
	}
}
