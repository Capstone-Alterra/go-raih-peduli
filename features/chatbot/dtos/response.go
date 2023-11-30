package dtos

type ResChatReply struct {
	Question string `json:"question"`
	Reply    string `json:"reply"`
}

type ResNewsContent struct {
	Content string `json:"content"`
}