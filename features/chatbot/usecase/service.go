package usecase

import (
	"raihpeduli/features/chatbot"
	"raihpeduli/features/chatbot/dtos"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
)

type service struct {
	model chatbot.Repository
}

func New(model chatbot.Repository) chatbot.Usecase {
	return &service {
		model: model,
	}
}

func (svc *service) FindAll(page, size int) []dtos.ResChatbot {
	var chatbots []dtos.ResChatbot

	chatbotsEnt := svc.model.Paginate(page, size)

	for _, chatbot := range chatbotsEnt {
		var data dtos.ResChatbot

		if err := smapping.FillStruct(&data, smapping.MapFields(chatbot)); err != nil {
			log.Error(err.Error())
		} 
		
		chatbots = append(chatbots, data)
	}

	return chatbots
}

func (svc *service) FindByID(chatbotID int) *dtos.ResChatbot {
	res := dtos.ResChatbot{}
	chatbot := svc.model.SelectByID(chatbotID)

	if chatbot == nil {
		return nil
	}

	err := smapping.FillStruct(&res, smapping.MapFields(chatbot))
	if err != nil {
		log.Error(err)
		return nil
	}

	return &res
}

func (svc *service) Create(newChatbot dtos.InputChatbot) *dtos.ResChatbot {
	chatbot := chatbot.Chatbot{}
	
	err := smapping.FillStruct(&chatbot, smapping.MapFields(newChatbot))
	if err != nil {
		log.Error(err)
		return nil
	}

	chatbotID := svc.model.Insert(chatbot)

	if chatbotID == -1 {
		return nil
	}

	resChatbot := dtos.ResChatbot{}
	errRes := smapping.FillStruct(&resChatbot, smapping.MapFields(newChatbot))
	if errRes != nil {
		log.Error(errRes)
		return nil
	}

	return &resChatbot
}

func (svc *service) Modify(chatbotData dtos.InputChatbot, chatbotID int) bool {
	newChatbot := chatbot.Chatbot{}

	err := smapping.FillStruct(&newChatbot, smapping.MapFields(chatbotData))
	if err != nil {
		log.Error(err)
		return false
	}

	newChatbot.ID = chatbotID
	rowsAffected := svc.model.Update(newChatbot)

	if rowsAffected <= 0 {
		log.Error("There is No Chatbot Updated!")
		return false
	}
	
	return true
}

func (svc *service) Remove(chatbotID int) bool {
	rowsAffected := svc.model.DeleteByID(chatbotID)

	if rowsAffected <= 0 {
		log.Error("There is No Chatbot Deleted!")
		return false
	}

	return true
}