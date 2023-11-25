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
	chatbots.POST("", handler.SendMessage(), m.AuthorizeJWT(jwt, -1, config.SECRET))
	
	chatbots.DELETE("/:user_id", handler.DeleteChatHistory(), m.AuthorizeJWT(jwt, 2, config.SECRET))
}