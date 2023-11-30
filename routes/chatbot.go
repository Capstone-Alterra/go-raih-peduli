package routes

import (
	"raihpeduli/config"
	"raihpeduli/features/chatbot"
	"raihpeduli/helpers"
	m "raihpeduli/middlewares"

	"github.com/labstack/echo/v4"
)

func Chatbots(e *echo.Echo, handler chatbot.Handler, jwt helpers.JWTInterface, config config.ProgramConfig) {
	chatbots := e.Group("/chatbots")

	chatbots.GET("", handler.GetChatHistory(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	chatbots.POST("", handler.SendQuestion(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	
	chatbots.DELETE("", handler.DeleteChatHistory(), m.AuthorizeJWT(jwt, 1, config.SECRET))

	e.POST("/generate-content", handler.GetNewsContentGeneration(), m.AuthorizeJWT(jwt, 2, config.SECRET))
}