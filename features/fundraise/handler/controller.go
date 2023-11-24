package handler

import (
	"mime/multipart"
	"raihpeduli/helpers"
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

func (ctl *controller) GetFundraises(suffix string) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)

		searchAndFilter := dtos.SearchAndFilter{}
		ctx.Bind(&searchAndFilter)

		userID := 0

		if ctx.Get("user_id") != nil {
			userID = ctx.Get("user_id").(int)
		}

		fundraises, totalData := ctl.service.FindAll(pagination, searchAndFilter, userID, suffix)

		if fundraises == nil {
			return ctx.JSON(404, helper.Response("fundraises not found"))
		}

		page := pagination.Page
		pageSize := pagination.PageSize

		if page <= 0 || pageSize <= 0 {
			page = 1
			pageSize = 10
		}

		paginationResponse := helpers.PaginationResponse(page, pageSize, int(totalData))

		return ctx.JSON(200, helper.Response("success", map[string]any{
			"data":       fundraises,
			"pagination": paginationResponse,
		}))
	}
}

func (ctl *controller) FundraiseDetails() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		fundraiseID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		userID := 0

		if ctx.Get("user_id") != nil {
			userID = ctx.Get("user_id").(int)
		}

		fundraise := ctl.service.FindByID(fundraiseID, userID)

		if fundraise == nil {
			return ctx.JSON(404, helper.Response("fundraise not found"))
		}

		return ctx.JSON(200, helper.Response("success", map[string]any{
			"data": fundraise,
		}))
	}
}

func (ctl *controller) CreateFundraise() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputFundraise{}

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

		fundraise, errMap, err := ctl.service.Create(input, userID.(int), file)

		if errMap != nil {
			return ctx.JSON(400, helper.Response("error missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		return ctx.JSON(200, helper.Response("success created fundraise", map[string]any{
			"data": fundraise,
		}))
	}
}

func (ctl *controller) UpdateFundraise() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputFundraise{}

		fundraiseID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		fundraise := ctl.service.FindByID(fundraiseID, 0)

		if fundraise == nil {
			return ctx.JSON(404, helper.Response("fundraise not found"))
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

		errMap, err := ctl.service.Modify(input, file, *fundraise)

		if errMap != nil {
			return ctx.JSON(400, helper.Response("error missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		return ctx.JSON(200, helper.Response("success updated fundraise"))
	}
}

func (ctl *controller) UpdateFundraiseStatus() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputFundraiseStatus{}

		fundraiseID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		fundraise := ctl.service.FindByID(fundraiseID, 0)

		if fundraise == nil {
			return ctx.JSON(404, helper.Response("fundraise not found"))
		}

		ctx.Bind(&input)

		errMap, err := ctl.service.ModifyStatus(input, *fundraise)

		if errMap != nil {
			return ctx.JSON(400, helper.Response("error missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(500, helper.Response(err.Error()))
		}

		return ctx.JSON(200, helper.Response("success updated fundraise status"))
	}
}

func (ctl *controller) DeleteFundraise() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		fundraiseID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		fundraise := ctl.service.FindByID(fundraiseID, 0)

		if fundraise == nil {
			return ctx.JSON(404, helper.Response("fundraise not found"))
		}

		delete := ctl.service.Remove(fundraiseID)

		if !delete {
			return ctx.JSON(500, helper.Response("something went wrong"))
		}

		return ctx.JSON(200, helper.Response("fundraise success deleted"))
	}
}
