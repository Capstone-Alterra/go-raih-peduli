package routes

import (
	"raihpeduli/features/customer"

	"github.com/labstack/echo/v4"
)

func Customers(e *echo.Echo, handler customer.Handler) {
	customers := e.Group("/customers")

	customers.GET("", handler.GetCustomers())
	customers.POST("", handler.CreateCustomer())

	customers.GET("/:id", handler.CustomerDetails())
	customers.PUT("/:id", handler.UpdateCustomer())
	customers.DELETE("/:id", handler.DeleteCustomer())
}
