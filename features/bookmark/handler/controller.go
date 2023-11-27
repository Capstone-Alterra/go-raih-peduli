package handler

import (
	"raihpeduli/helpers"

	"raihpeduli/features/bookmark"
	"raihpeduli/features/bookmark/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service bookmark.Usecase
}

func New(service bookmark.Usecase) bookmark.Handler {
	return &controller{
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetBookmarksByUserID() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)

		size := pagination.PageSize

		if size <= 0 {
			size = 10
		}

		userID := ctx.Get("user_id").(int)

		bookmarks := ctl.service.FindAll(size, userID)

		if bookmarks == nil {
			return ctx.JSON(404, helpers.Response("there is no bookmarks"))
		}

		return ctx.JSON(200, helpers.Response("success", map[string]any{
			"data": bookmarks,
		}))
	}
}

func (ctl *controller) BookmarkAPost() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputBookmarkPost{}

		ctx.Bind(&input)

		userID := ctx.Get("user_id").(int)

		_, errMap, err := ctl.service.SetBookmark(input, userID)

		if errMap != nil {
			return ctx.JSON(400, helpers.Response("error missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("success bookmarked a post"))
	}
}

func (ctl *controller) UnBookmarkAPost() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		bookmarkID := ctx.Param("id")

		userID := ctx.Get("user_id").(int)

		bookmark := ctl.service.FindByID(bookmarkID)

		if bookmark == nil {
			return ctx.JSON(404, helpers.Response("post not found"))
		}

		delete, err := ctl.service.UnsetBookmark(bookmarkID, bookmark, userID)

		if !delete {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("bookmark success deleted", nil))
	}
}
