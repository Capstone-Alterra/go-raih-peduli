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
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)

		if pagination.Page < 1 || pagination.Size < 1 {
			pagination.Page = 1
			pagination.Size = 20
		}

		page := pagination.Page
		size := pagination.Size
		users, totalData := ctl.service.FindAll(page, size)

		if users == nil {
			return ctx.JSON(404, helpers.Response("There is No Users!"))
		}

		paginationResponse := helpers.PaginationResponse(page, size, int(totalData))

		return ctx.JSON(200, helpers.Response("Success!", map[string]any{
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
			return ctx.JSON(404, helpers.Response("User Not Found!"))
		}

		return ctx.JSON(200, helpers.Response("Success!", map[string]any{
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
			return ctx.JSON(400, helpers.Response("Bad Request!", map[string]any{
				"error": err.Error(),
			}))
		}

		return ctx.JSON(200, helpers.Response("Success!", map[string]any{
			"data": user,
		}))
	}
}

func (ctl *controller) UpdateUser() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputUpdate{}

		userID := ctx.Get("user_id").(int)

		user := ctl.service.FindByID(userID)

		if user == nil {
			return ctx.JSON(404, helpers.Response("User Not Found!"))
		}

		ctx.Bind(&input)

		fileHeader, err := ctx.FormFile("profile_picture")
		var file multipart.File

		if err == nil {
			formFile, err := fileHeader.Open()

			if err != nil {
				return ctx.JSON(500, helpers.Response("something went wrong"))
			}

			file = formFile
		}

		update, errMap := ctl.service.Modify(input, file, *user)
		if errMap != nil {
			return ctx.JSON(400, helpers.Response("missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if !update {
			return ctx.JSON(500, helpers.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helpers.Response("User Success Updated!"))
	}
}

func (ctl *controller) UpdateProfilePicture() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputUpdateProfilePicture{}

		userID := ctx.Get("user_id").(int)

		user := ctl.service.FindByID(userID)

		if user == nil {
			return ctx.JSON(404, helpers.Response("User Not Found!"))
		}

		fileHeader, err := ctx.FormFile("profile_picture")

		if err == nil {
			formFile, err := fileHeader.Open()

			if err != nil {
				return ctx.JSON(500, helpers.Response("something went wrong"))
			}

			input.ProfilePicture = formFile
		}

		update, errMap := ctl.service.ModifyProfilePicture(input, *user)
		if errMap != nil {
			return ctx.JSON(400, helpers.Response("missing some data", map[string]any{
				"error": errMap,
			}))
		}

		if !update {
			return ctx.JSON(500, helpers.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helpers.Response("User Success Updated!"))
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
			return ctx.JSON(404, helpers.Response("User Not Found!"))
		}

		delete := ctl.service.Remove(userID)

		if !delete {
			return ctx.JSON(500, helpers.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helpers.Response("User Success Deleted!", nil))
	}
}

func (ctl *controller) VerifyEmail() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.VerifyOTP{}

		ctx.Bind(&input)

		verifyOTP := ctl.service.ValidateVerification(input.OTP)
		if !verifyOTP {
			return ctx.JSON(400, helpers.Response("Incorrect / Expired OTP"))
		}

		return ctx.JSON(200, helpers.Response("Success verify email!"))
	}
}

func (ctl *controller) ForgetPassword() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var email dtos.ForgetPassword

		ctx.Bind(&email)

		err := ctl.service.ForgetPassword(email)
		if err != nil {
			return ctx.JSON(404, helpers.Response("User Not Found!"))
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
			return ctx.JSON(400, helpers.Response("Incorrect / Expired OTP"))
		}

		return ctx.JSON(200, helpers.Response("Success verify email!", map[string]any{
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
			return ctx.JSON(500, helpers.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helpers.Response("Success Reset Password!"))
	}
}

func (ctl *controller) MyProfile() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		userID := ctx.Get("user_id").(int)

		user := ctl.service.FindByID(userID)
		if user == nil {
			return ctx.JSON(404, helpers.Response("User Not Found!"))
		}

		return ctx.JSON(200, helpers.Response("Success!", map[string]any{
			"data": user,
		}))
	}
}

func (ctl *controller) RefreshJWT() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		jwt := dtos.RefreshJWT{}
		ctx.Bind(&jwt)

		refershJWT, err := ctl.service.RefreshJWT(jwt)
		if err != nil {
			if err.Error() == "validate token failed" {
				return ctx.JSON(400, helpers.Response("invalid jwt token"))
			}

			return ctx.JSON(500, helpers.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helpers.Response("Success!", map[string]any{
			"data": refershJWT,
		}))
	}
}
