package handler

import (
	"mime/multipart"
	"raihpeduli/helpers"
	"strconv"

	"raihpeduli/features/user"
	"raihpeduli/features/user/dtos"

	"github.com/labstack/echo/v4"
)

type controller struct {
	service user.Usecase
}

func New(service user.Usecase) user.Handler {
	return &controller{
		service: service,
	}
}

func (ctl *controller) GetUsers() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		searchAndFilter := dtos.SearchAndFilter{}
		ctx.Bind(&searchAndFilter)

		if searchAndFilter.Page < 1 || searchAndFilter.PageSize < 1 {
			searchAndFilter.Page = 1
			searchAndFilter.PageSize = 20
		}

		users, totalData := ctl.service.FindAll(searchAndFilter)

		if users == nil {
			return ctx.JSON(404, helpers.Response("there is no users"))
		}

		paginationResponse := helpers.PaginationResponse(searchAndFilter.Page, searchAndFilter.PageSize, int(totalData))

		return ctx.JSON(200, helpers.Response("success", map[string]any{
			"data":       users,
			"pagination": paginationResponse,
		}))
	}
}

func (ctl *controller) UserDetails() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helpers.Response(err.Error()))
		}

		user := ctl.service.FindByID(userID)

		if user == nil {
			return ctx.JSON(404, helpers.Response("user not found"))
		}

		return ctx.JSON(200, helpers.Response("success", map[string]any{
			"data": user,
		}))
	}
}

func (ctl *controller) CreateUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputUser{}

		ctx.Bind(&input)

		user, errMap, err := ctl.service.Create(input)
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

func (ctl *controller) UpdateUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var roleID = ctx.Get("role_id")
		var userID int
		var input dtos.InputUpdate

		if roleID == 2 {
			userID, _ = strconv.Atoi(ctx.Param("id"))
		} else {
			userID = ctx.Get("user_id").(int)
		}

		user := ctl.service.FindByID(userID)
		if user == nil {
			return ctx.JSON(404, helpers.Response("user not found"))
		}

		ctx.Bind(&input)

		var file multipart.File
		fileHeader, err := ctx.FormFile("profile_picture")

		if err == nil {
			formFile, err := fileHeader.Open()

			if err != nil {
				return ctx.JSON(500, helpers.Response("something went wrong"))
			}

			file = formFile
		}

		err, errMap := ctl.service.Modify(input, file, *user)
		if errMap != nil {
			return ctx.JSON(400, helpers.Response("missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("success updated user"))
	}
}

func (ctl *controller) UpdateProfilePicture() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputUpdateProfilePicture{}

		userID := ctx.Get("user_id").(int)

		user := ctl.service.FindByID(userID)

		if user == nil {
			return ctx.JSON(404, helpers.Response("user not found"))
		}

		fileHeader, err := ctx.FormFile("profile_picture")

		if err == nil {
			formFile, err := fileHeader.Open()

			if err != nil {
				return ctx.JSON(500, helpers.Response("something went wrong"))
			}

			input.ProfilePicture = formFile
		}

		err, errMap := ctl.service.ModifyProfilePicture(input, *user)
		if errMap != nil {
			return ctx.JSON(400, helpers.Response("missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("success updated user"))
	}
}

func (ctl *controller) DeleteUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helpers.Response(err.Error()))
		}

		user := ctl.service.FindByID(userID)

		if user == nil {
			return ctx.JSON(404, helpers.Response("user not found"))
		}

		err = ctl.service.Remove(userID)
		if err != nil {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("success deleted user", nil))
	}
}

func (ctl *controller) VerifyEmail() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.VerifyOTP{}

		ctx.Bind(&input)

		verifyOTP := ctl.service.ValidateVerification(input.OTP)
		if !verifyOTP {
			return ctx.JSON(400, helpers.Response("incorrect / expired OTP"))
		}

		return ctx.JSON(200, helpers.Response("success verify email"))
	}
}

func (ctl *controller) ForgetPassword() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var email dtos.ForgetPassword

		ctx.Bind(&email)

		err := ctl.service.ForgetPassword(email)
		if err != nil {
			return ctx.JSON(404, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("OTP has been sent via email"))
	}
}

func (ctl *controller) VerifyOTP() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var input dtos.VerifyOTP

		ctx.Bind(&input)

		token := ctl.service.VerifyOTP(input.OTP)
		if token == "" {
			return ctx.JSON(400, helpers.Response("incorrect / expired OTP"))
		}

		return ctx.JSON(200, helpers.Response("success verify email", map[string]any{
			"access_token": token,
		}))
	}
}

func (ctl *controller) ResetPassword() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var input dtos.ResetPassword

		ctx.Bind(&input)

		err := ctl.service.ResetPassword(input)

		if err != nil {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("success reset password"))
	}
}

func (ctl *controller) MyProfile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get("user_id").(int)

		user := ctl.service.MyProfile(userID)
		if user == nil {
			return ctx.JSON(404, helpers.Response("user not found"))
		}

		return ctx.JSON(200, helpers.Response("success", map[string]any{
			"data": user,
		}))
	}
}

func (ctl *controller) CheckPassword() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get("user_id").(int)

		var input dtos.CheckPassword

		ctx.Bind(&input)

		errMap, err := ctl.service.CheckPassword(input, userID)
		if errMap != nil {
			return ctx.JSON(400, helpers.Response("missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("success check password"))
	}
}

func (ctl *controller) ChangePassword() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get("user_id").(int)

		var input dtos.ChangePassword

		ctx.Bind(&input)

		errMap, err := ctl.service.ChangePassword(input, userID)
		if errMap != nil {
			return ctx.JSON(400, helpers.Response("missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if err != nil {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("success change password"))
	}
}

func (ctl *controller) AddPersonalization() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get("user_id").(int)

		var input dtos.InputPersonalization

		ctx.Bind(&input)

		err := ctl.service.AddPersonalization(userID, input)
		if err != nil {
			return ctx.JSON(500, helpers.Response(err.Error()))
		}

		return ctx.JSON(200, helpers.Response("success add personalization"))
	}
}
