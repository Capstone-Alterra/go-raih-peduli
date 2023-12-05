package chatbot

import "time"

type ChatHistory struct {
	UserID           int                `bson:"user_id"`
	QuestionAndReply []QuestionAndReply `bson:"question_reply"`
}

type QuestionAndReply struct {
	Question     string    `bson:"question"`
	QuestionTime time.Time `bson:"question_time"`
	Reply        string    `bson:"reply"`
	ReplyTime 	 time.Time `bson:"reply_time"`
}

type User struct {
	ID    int
	Email string
}

type QuestionAndPrompt struct {
	Question string `json:"question"`
	Prompt   string `json:"prompt"`
}