package chatbot

import (
	"raihpeduli/features/chatbot/dtos"

	"github.com/labstack/echo/v4"
)

type Repository interface {
	Paginate(page, size int) []Chatbot
	Insert(newChatbot Chatbot) int64
	SelectByID(chatbotID int) *Chatbot
	Update(chatbot Chatbot) int64
	DeleteByID(chatbotID int) int64
}

type Usecase interface {
	FindAll(page, size int) []dtos.ResChatbot
	FindByID(chatbotID int) *dtos.ResChatbot
	Create(newChatbot dtos.InputMessage) *dtos.ResChatbot
	Modify(chatbotData dtos.InputMessage, chatbotID int) bool
	Remove(chatbotID int) bool
}

type Handler interface {
	GetChatHistory() echo.HandlerFunc
	SendMessage() echo.HandlerFunc
	DeleteChatHistory() echo.HandlerFunc
}
