package handler

import (
	"encoding/json"
	"raihpeduli/helpers"
	helper "raihpeduli/helpers"
	"strconv"

	"raihpeduli/features/transaction"
	"raihpeduli/features/transaction/dtos"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type controller struct {
	service transaction.Usecase
}

func New(service transaction.Usecase) transaction.Handler {
	return &controller{
		service: service,
	}
}

var validate *validator.Validate

func (ctl *controller) GetTransactions() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		pagination := dtos.Pagination{}
		var roleID = ctx.Get("role_id")

		ctx.Bind(&pagination)

		keyword := ctx.QueryParam("fullname")
		userID := ctx.Get("user_id").(int)

		page := pagination.Page
		size := pagination.PageSize

		if page <= 0 || size <= 0 {
			return ctx.JSON(400, helper.Response("Please provide query `page` and `size` in number!"))
		}

		transactions, totalData := ctl.service.FindAll(page, size, roleID.(int), userID, keyword)

		if transactions == nil {
			return ctx.JSON(404, helper.Response("There is No Transactions!"))
		}

		paginationResponse := helpers.PaginationResponse(page, size, int(totalData))

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data":       transactions,
			"pagination": paginationResponse,
		}))
	}
}

func (ctl *controller) TransactionDetails() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		transactionID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		transaction := ctl.service.FindByID(transactionID)

		if transaction == nil {
			return ctx.JSON(404, helper.Response("Transaction Not Found!"))
		}

		return ctx.JSON(200, helper.Response("Success!", map[string]any{
			"data": transaction,
		}))
	}
}

func (ctl *controller) CreateTransaction() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputTransaction{}

		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())

		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any{
				"error": errMap,
			}))
		}

		userID := ctx.Get("user_id").(int)
		transaction, err, validate := ctl.service.Create(userID, input)

		if err != nil {
			return ctx.JSON(400, helper.BuildErrorResponse(err.Error()))
		}

		if transaction == nil {
			return ctx.JSON(500, helper.Response("Something went Wrong!", nil))
		}

		if validate != nil {
			return ctx.JSON(400, helper.Response("error missing some data", map[string]any{
				"error": validate,
			}))
		}

		return ctx.JSON(200, map[string]interface{}{
			"message": "Success!",
			"data":    transaction,
		})
	}
}

func (ctl *controller) UpdateTransaction() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		input := dtos.InputTransaction{}

		transactionID, errParam := strconv.Atoi(ctx.Param("id"))

		if errParam != nil {
			return ctx.JSON(400, helper.Response(errParam.Error()))
		}

		transaction := ctl.service.FindByID(transactionID)

		if transaction == nil {
			return ctx.JSON(404, helper.Response("Transaction Not Found!"))
		}

		ctx.Bind(&input)

		validate = validator.New(validator.WithRequiredStructEnabled())
		err := validate.Struct(input)

		if err != nil {
			errMap := helpers.ErrorMapValidation(err)
			return ctx.JSON(400, helper.Response("Bad Request!", map[string]any{
				"error": errMap,
			}))
		}

		update := ctl.service.Modify(input, transactionID)

		if !update {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Transaction Success Updated!"))
	}
}

func (ctl *controller) DeleteTransaction() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		transactionID, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			return ctx.JSON(400, helper.Response(err.Error()))
		}

		transaction := ctl.service.FindByID(transactionID)

		if transaction == nil {
			return ctx.JSON(404, helper.Response("Transaction Not Found!"))
		}

		delete := ctl.service.Remove(transactionID)

		if !delete {
			return ctx.JSON(500, helper.Response("Something Went Wrong!"))
		}

		return ctx.JSON(200, helper.Response("Transaction Success Deleted!", nil))
	}
}

func (ctl *controller) Notifications() echo.HandlerFunc {
	return func(c echo.Context) error {
		var notificationPayload map[string]any

		if err := json.NewDecoder(c.Request().Body).Decode(&notificationPayload); err != nil {
			return c.JSON(500, helper.Response("Something Went Wrong!"))
		}

		err := ctl.service.Notifications(notificationPayload)
		if err != nil {
			return c.JSON(505, helper.Response("Something Went Wrong!"))
		}

		return c.JSON(200, echo.Map{"status": "ok"})
	}
}
