package handler

import (
	"raihpeduli/helpers"
	helper "raihpeduli/helpers"
	"strconv"

	"raihpeduli/features/customer"
	"raihpeduli/features/customer/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service customer.Usecase
}

func New(service customer.Usecase) customer.Handler {
	return &controller {
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetCustomers() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		pagination := dtos.Pagination{}
		ctx.Bind(&pagination)
		
		page := pagination.Page
		size := pagination.Size

		if page <= 0 || size <= 0 {
			return ctx.JSON(400, helper.Response("Please provide query `page` and `size` in number!"))
		}

		customers := ctl.service.FindAll(page, size)

		if customers == nil {
			return ctx.JSON(404, helper.Response("There is No Customers!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": customers,
		}))
	}
}


func (ctl *controller) CustomerDetails() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		customerID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		customer := ctl.service.FindByID(customerID)

		if customer == nil {
			return ctx.JSON(404, helper.Response("Customer Not Found!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": customer,
		}))
	}
}

func (ctl *controller) CreateCustomer() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		input := dtos.InputCustomer{}

		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())

		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any {
				"error": errMap,
			}))
		}

		customer := ctl.service.Create(input)

		if customer == nil {
			return ctx.JSON(500, helper.Response("Something went Wrong!", nil))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any {
			"data": customer,
		}))
	}
}

func (ctl *controller) UpdateCustomer() echo.HandlerFunc {
	return func (ctx echo.Context) error {
		input := dtos.InputCustomer{}

		customerID, errParam := strconv.Atoi(ctx.Param("id"))

		if errParam != nil {
			return ctx.JSON(400, helper.Response(errParam.Error()))
		}

		customer := ctl.service.FindByID(customerID)

		if customer == nil {
			return ctx.JSON(404, helper.Response("Customer Not Found!"))
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

		update := ctl.service.Modify(input, customerID)

		if !update {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Customer Success Updated!"))
	}
}

func (ctl *controller) DeleteCustomer() echo.HandlerFunc {
	return func (ctx echo.Context) error  {
		customerID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		customer := ctl.service.FindByID(customerID)

		if customer == nil {
			return ctx.JSON(404, helper.Response("Customer Not Found!"))
		}

		delete := ctl.service.Remove(customerID)

		if !delete {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Customer Success Deleted!", nil))
	}
}
