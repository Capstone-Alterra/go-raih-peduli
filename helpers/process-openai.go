package helpers

import (
	"context"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

type openAI struct {
	API_KEY string
}

func NewOpenAI(API_KEY string) OpenAIInterface {
	return &openAI{
		API_KEY: API_KEY,
	}
}

func (oai *openAI) GetAppInformation(question string, qnaList map[string]string) (string, error) {
	var listOfQuestionNPrompt string

	for question, prompt := range qnaList {
		listOfQuestionNPrompt += fmt.Sprintf("question : '%s'\nprompt : '%s'\n----------------------------------------------------------\n", question, prompt)
	}

	prompt := fmt.Sprintf("Given the following list of question of an APP QNA information: \n\n%s\n\nAnd given the user's question: '%s'\n\nPlease provide a detailed answer based on the information in the QNA list. Answer in bahasa if the question is in bahasa. If there is no match to the question given, provide this message 'there is no information related to that question in this app' as in bahasa.", listOfQuestionNPrompt, question)

	client := openai.NewClient(oai.API_KEY)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}
	
	return resp.Choices[0].Message.Content, nil
}

func (oai *openAI) GetNewsContent(prompt string) (string, error) {
	fmt.Println(prompt)
	client := openai.NewClient(oai.API_KEY)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			// MaxTokens: 200,
			// Stop: []string{"."},
		},
	)

	if err != nil {
		return "", err
	}
	
	return resp.Choices[0].Message.Content, nil
}