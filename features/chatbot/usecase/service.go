package usecase

import (
	"errors"
	"raihpeduli/features/chatbot"
	"raihpeduli/features/chatbot/dtos"
	"raihpeduli/helpers"

	"github.com/labstack/gommon/log"
	"github.com/mashingan/smapping"
	"github.com/sirupsen/logrus"
)

type service struct {
	model chatbot.Repository
	validation helpers.ValidationInterface
	openAI helpers.OpenAIInterface
}

func New(model chatbot.Repository, validation helpers.ValidationInterface, openAI helpers.OpenAIInterface) chatbot.Usecase {
	return &service {
		model: model,
		validation: validation,
		openAI: openAI,
	}
}

func (svc *service) FindAllChat(userID int) []dtos.ResChatReply {
	var res []dtos.ResChatReply

	chatHistories, err := svc.model.SelectByUserID(userID)

	if err != nil {
		logrus.Error(err)
		return nil
	}

	for _, chatbot := range chatHistories.QuestionAndReply {
		var data dtos.ResChatReply

		if err := smapping.FillStruct(&data, smapping.MapFields(chatbot)); err != nil {
			log.Error(err.Error())
		} 
		
		res = append(res, data)
	}

	return res
}

func (svc *service) SetReplyMessage(input dtos.InputMessage, userID int) (*dtos.ResChatReply, []string, error) {
	if errMap := svc.validation.ValidateRequest(input); errMap != nil {
		return nil, errMap, errors.New("message must not be empty") 
	}

	// _, err := svc.model.ReadQuestionNPrompts()
	
	reply, err := svc.openAI.GetReplyFromGPT(input.Message, map[string]string{})

	if err != nil {
		return nil, nil, err
	}

	var chatMessage = chatbot.QuestionAndReply{
		Question: input.Message,
		Reply: reply,
	}
	
	if userID != 0 {
		if err := svc.model.SaveChat(chatMessage, userID); err != nil {
			return nil, nil, err
		}
	}

	var res = dtos.ResChatReply{
		Question: input.Message,
		Reply: reply,
	}

	return &res, nil, nil
}

func (svc *service) ClearHistory(userID int) error {
	if err := svc.model.DeleteByUserID(userID); err != nil {
		return err
	}

	return nil
}