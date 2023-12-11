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

func (svc *service) SetContentForNews(input dtos.InputMessage) (*dtos.ResNewsContent, []string, error) {
	if errMap := svc.validation.ValidateRequest(input); errMap != nil {
		return nil, errMap, errors.New("message must not be empty") 
	}

	reply, err := svc.openAI.GetNewsContent(input.Message)

	if err != nil {
		return nil, nil, err
	}

	return &dtos.ResNewsContent{
		Content: reply,
	}, nil, nil
}

func (svc *service) SetReplyMessage(input dtos.InputMessage, userID int) (*dtos.ResChatReply, []string, error) {
	if errMap := svc.validation.ValidateRequest(input); errMap != nil {
		return nil, errMap, errors.New("message must not be empty") 
	}

	data, err := svc.model.ReadQuestionNPrompts()
	if err != nil {
		return nil, nil, err
	}

	var chatMessage = chatbot.QuestionAndReply{
		Question: input.Message,
		QuestionTime: svc.model.GetTimeNow(),
	}

	reply, err := svc.openAI.GetAppInformation(input.Message, data)

	if err != nil {
		return nil, nil, err
	}

	chatMessage.Reply = reply
	chatMessage.ReplyTime = svc.model.GetTimeNow()
	
	if userID != 0 {
		if err := svc.model.SaveChat(chatMessage, userID); err != nil {
			return nil, nil, err
		}
	}

	return &dtos.ResChatReply{
		Question: input.Message,
		Reply: reply,
		QuestionTime: chatMessage.QuestionTime,
		ReplyTime: chatMessage.ReplyTime,
	}, nil, nil
}

func (svc *service) ClearHistory(userID int) error {
	if err := svc.model.DeleteByUserID(userID); err != nil {
		return err
	}

	return nil
}