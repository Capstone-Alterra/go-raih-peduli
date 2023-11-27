package usecase

import (
	"errors"
	"raihpeduli/features/chatbot"
	"raihpeduli/features/chatbot/dtos"
	"raihpeduli/features/chatbot/mocks"
	helperMocks "raihpeduli/helpers/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindAllChat(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var openai = helperMocks.NewOpenAIInterface(t)
	var service = New(repository, validation, openai)

	var chatHistories = chatbot.ChatHistory{
		UserID: 1,
		QuestionAndReply: []chatbot.QuestionAndReply{
			{
				Question: "",
				Reply: "",
			},
		},
	}

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("SelectByUserID", userID).Return(&chatHistories, nil).Once()

		res := service.FindAllChat(userID)
		assert.Equal(t, chatHistories.QuestionAndReply[0].Question, res[0].Question)
		repository.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repository.On("SelectByUserID", userID).Return(nil, errors.New("record not found")).Once()

		res := service.FindAllChat(userID)
		assert.Nil(t, res)
		repository.AssertExpectations(t)
	})
}

func TestSetReplyMessage(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var openai = helperMocks.NewOpenAIInterface(t)
	var service = New(repository, validation, openai)

	var input = dtos.InputMessage{
		Message: "bagaimana cara donasi pada program penggalangan dana?",
	}

	var data = map[string]string {
		"question": "prompt",
	}

	var chatMessage = chatbot.QuestionAndReply{
		Question: input.Message,
		Reply: "this answer generated by gpt",
	}

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("ReadQuestionNPrompts").Return(data, nil).Once()
		openai.On("GetReplyFromGPT", input.Message, data).Return("this answer generated by gpt", nil).Once()
		repository.On("SaveChat", chatMessage, userID).Return(nil).Once()

		res, errMap, err := service.SetReplyMessage(input, userID)
		assert.Equal(t, res.Question, input.Message)
		assert.Nil(t, errMap)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error When Save Chat", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("ReadQuestionNPrompts").Return(data, nil).Once()
		openai.On("GetReplyFromGPT", input.Message, data).Return("this answer generated by gpt", nil).Once()
		repository.On("SaveChat", chatMessage, userID).Return(errors.New("can't insert new chat")).Once()

		res, errMap, err := service.SetReplyMessage(input, userID)
		assert.Nil(t, res)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error Request To OpenAI", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("ReadQuestionNPrompts").Return(data, nil).Once()
		openai.On("GetReplyFromGPT", input.Message, data).Return("", errors.New("error when request to openai")).Once()

		res, errMap, err := service.SetReplyMessage(input, userID)
		assert.Nil(t, res)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed : Error Read File JSON", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return(nil).Once()
		repository.On("ReadQuestionNPrompts").Return(nil, errors.New("can't read file")).Once()

		res, errMap, err := service.SetReplyMessage(input, userID)
		assert.Nil(t, res)
		assert.Nil(t, errMap)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
	
	t.Run("Failed : Error Validation Request", func(t *testing.T) {
		validation.On("ValidateRequest", input).Return([]string{"message is required"}).Once()

		res, errMap, err := service.SetReplyMessage(input, userID)
		assert.Nil(t, res)
		assert.NotNil(t, errMap)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
}

func TestClearHistory(t *testing.T) {
	var repository = mocks.NewRepository(t)
	var validation = helperMocks.NewValidationInterface(t)
	var openai = helperMocks.NewOpenAIInterface(t)
	var service = New(repository, validation, openai)

	var userID = 1

	t.Run("Success", func(t *testing.T) {
		repository.On("DeleteByUserID", userID).Return(nil).Once()

		err := service.ClearHistory(userID)
		assert.Nil(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repository.On("DeleteByUserID", userID).Return(errors.New("error when delete")).Once()

		err := service.ClearHistory(userID)
		assert.NotNil(t, err)
		repository.AssertExpectations(t)
	})
}