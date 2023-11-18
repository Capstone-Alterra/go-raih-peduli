package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/transaction"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func Transactions(e *echo.Echo, handler transaction.Handler, jwt helpers.JWTInterface, config config.ProgramConfig) {
	transactions := e.Group("/transactions")

	transactions.GET("", handler.GetTransactions(), m.AuthorizeJWT(jwt, 0, config.SECRET))
	transactions.POST("", handler.CreateTransaction(), m.AuthorizeJWT(jwt, 0, config.SECRET))
	transactions.POST("/notifications", handler.Notifications())

	transactions.GET("/:id", handler.TransactionDetails())
	transactions.PUT("/:id", handler.UpdateTransaction())
	transactions.DELETE("/:id", handler.DeleteTransaction())
}
