package repository

import (
	"context"
	"encoding/json"
	"os"
	"raihpeduli/features/chatbot"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
	collection *mongo.Collection
}

func New(db *gorm.DB, collection *mongo.Collection) chatbot.Repository {
	return &model {
		db: db,
		collection: collection,
	}
}

func (mdl *model) SaveChat(questionNReply chatbot.QuestionAndReply, userID int) error {
	var user chatbot.User
	
	if err := mdl.collection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&user); err != nil {
		mdl.collection.InsertOne(context.Background(), chatbot.ChatHistory{
			UserID: userID,
			QuestionAndReply: []chatbot.QuestionAndReply{
				questionNReply,
			},
		})
	} else {
		if _, err := mdl.collection.UpdateOne(context.Background(), bson.M{"user_id": userID}, bson.M{"$push": bson.M{"question_reply": questionNReply}}); err != nil {
			return err
		}
	}

	return nil
}

func (mdl *model) SelectUserByID(userID int) (*chatbot.User, error) {
	var user chatbot.User

	if err := mdl.db.Table("users").Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (mdl *model) SelectByUserID(userID int) (*chatbot.ChatHistory, error) {
	
	var chatHistories chatbot.ChatHistory
	if err := mdl.collection.FindOne(context.Background(), bson.M{"user_id": userID}).Decode(&chatHistories); err != nil {
		return nil, err
	}

	return &chatHistories, nil
}

func (mdl *model) DeleteByUserID(userID int) error {
	if _, err := mdl.collection.DeleteOne(context.Background(), bson.M{"user_id": userID}); err != nil {
		return err
	}

	return nil
}

func (mdl *model) ReadQuestionNPrompts() (map[string]string, error) {
	var filepath = "./features/chatbot/q-and-prompt.json"
	
	jsonData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var QAPairs []chatbot.QuestionAndPrompt
	if err = json.Unmarshal(jsonData, &QAPairs); err != nil {
		return nil, err
	}

	var data = make(map[string]string)
	for _, qa := range QAPairs {
		data[qa.Question] = qa.Prompt
	}

	return data, nil
}