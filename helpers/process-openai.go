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

func (oai *openAI) GetReplyFromGPT(question string, qnaList map[string]string) (string, error) {
	prompt := fmt.Sprintf("Given the following list of question of an APP QNA information: \nquestion : 'how to view available fundraising programs?'\nprompt : 'explain how a fundraising programs can be viewed through an app. The available one is the one that has status accepted, not pending or rejected'\n--------------------------------------------------------------------------------\nquestion : 'i want to know about fundraising programs?'\nprompt : 'explain what is a fundraising programs.'\n--------------------------------------------------------------------------------\nquestion : 'how to donate to a fundraising program?'\nprompt : 'explain step by step how to donate to a fundraising program or post. choose one of the fundraising details, click donate button, choose the payment type, and then pay the bill. generate better answers.'\n\nAnd given the user's question: '%s'\n\nPlease provide a detailed answer based on the information in the QNA list. Answer in bahasa if the question is in bahasa. If there is no match to the question given, provide this message 'there is no information related to that question in this app' as in bahasa.", question)

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
