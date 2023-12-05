package dtos

import "time"

type ResChatReply struct {
	Question     string    `json:"question"`
	QuestionTime time.Time `json:"question_time"`
	Reply        string    `json:"reply"`
	ReplyTime    time.Time `json:"reply_time"`
}

type ResNewsContent struct {
	Content string `json:"content"`
}