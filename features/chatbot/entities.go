package chatbot

type ChatHistory struct {
	UserID           int                `bson:"user_id"`
	QuestionAndReply []QuestionAndReply `bson:"question_reply"`
}

type QuestionAndReply struct {
	Question string `bson:"question"`
	Reply    string `bson:"reply"`
}

type User struct {
	ID    int
	Email string
}