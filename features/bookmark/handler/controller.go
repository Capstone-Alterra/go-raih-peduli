package handler

import (
	"raihpeduli/helpers"
	"strconv"

	"raihpeduli/features/bookmark"
	"raihpeduli/features/bookmark/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service bookmark.Usecase
}

func New(service bookmark.Usecase) bookmark.Handler {
	return &controller {
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetBookmarksByUserID() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)
		
		page := pagination.Page
		size := pagination.Size

		if page <= 0 || size <= 0 {
			page = 1
			size = 10
		}

		bookmarks := ctl.service.FindAll(page, size)

		if bookmarks == nil {
			return ctx.JSON(404, helpers.Response("There is No Bookmarks!"))
		}

		return ctx.JSON(200, helpers.Response("Success!", map[string]any {
			"data": bookmarks,
		}))
	}
}

func (ctl *controller) BookmarkAPost() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		input := dtos.InputBookmarkPost{}

		ctx.Bind(&input)

		userID := ctx.Get("user_id")

		_, err := ctl.service.SetBookmark(input, userID.(int))

		if err != nil {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("success bookmark a post"))
	}
}

func (ctl *controller) UnBookmarkAPost() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		bookmarkID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helpers.Response(err.Error()))
		}

		delete := ctl.service.UnsetBookmark(bookmarkID)

		if !delete {
			return ctx.JSON(500, helpers.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helpers.Response("Bookmark Success Deleted!", nil))
	}
}
