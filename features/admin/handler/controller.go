package handler

import (
	"raihpeduli/helpers"
	helper "raihpeduli/helpers"
	"strconv"

	"raihpeduli/features/admin"
	"raihpeduli/features/admin/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service admin.Usecase
}

func New(service admin.Usecase) admin.Handler {
	return &controller {
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetAdmins() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)
		
		page := pagination.Page
		size := pagination.Size

		if page <= 0 || size <= 0 {
			return ctx.JSON(400, helper.Response("Please provide query `page` and `size` in number!"))
		}

		admins := ctl.service.FindAll(page, size)

		if admins == nil {
			return ctx.JSON(404, helper.Response("There is No Admins!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": admins,
		}))
	}
}


func (ctl *controller) AdminDetails() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		adminID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		admin := ctl.service.FindByID(adminID)

		if admin == nil {
			return ctx.JSON(404, helper.Response("Admin Not Found!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": admin,
		}))
	}
}

func (ctl *controller) CreateAdmin() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		input := dtos.InputAdmin{}

		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())

		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any {
				"error": errMap,
			}))
		}

		admin := ctl.service.Create(input)

		if admin == nil {
			return ctx.JSON(500, helper.Response("Something went Wrong!", nil))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": admin,
		}))
	}
}

func (ctl *controller) UpdateAdmin() echo.HandlerFunc {
	return func (ctx echo.Context) error {
		input := dtos.InputAdmin{}

		adminID, errParam := strconv.Atoi(ctx.Param("id"))

		if errParam != nil {
			return ctx.JSON(400, helper.Response(errParam.Error()))
		}

		admin := ctl.service.FindByID(adminID)

		if admin == nil {
			return ctx.JSON(404, helper.Response("Admin Not Found!"))
		}
		
		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any {
				"error": errMap,
			}))
		}

		update := ctl.service.Modify(input, adminID)

		if !update {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Admin Success Updated!"))
	}
}

func (ctl *controller) DeleteAdmin() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		adminID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		admin := ctl.service.FindByID(adminID)

		if admin == nil {
			return ctx.JSON(404, helper.Response("Admin Not Found!"))
		}

		delete := ctl.service.Remove(adminID)

		if !delete {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Admin Success Deleted!", nil))
	}
}

func (ctl *controller) LoginAdmin() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		loginData := dtos.LoginAdmin{}

		if err := ctx.Bind(&loginData); err != nil {
			return ctx.JSON(400, helper.Response("Invalid request body!"))
		}

		loginRes, err := ctl.service.Login(loginData.Email, loginData.Password)
		if err != nil {
			return ctx.JSON(401, helper.Response("Invalid credentials!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": loginRes,
		}))
	}
}