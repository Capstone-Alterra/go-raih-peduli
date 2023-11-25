package repository

import (
	"raihpeduli/features/chatbot"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type model struct {
	db *gorm.DB
}

func New(db *gorm.DB) chatbot.Repository {
	return &model {
		db: db,
	}
}

func (mdl *model) Paginate(page, size int) []chatbot.Chatbot {
	var chatbots []chatbot.Chatbot

	offset := (page - 1) * size

	result := mdl.db.Offset(offset).Limit(size).Find(&chatbots)
	
	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return chatbots
}

func (mdl *model) Insert(newChatbot chatbot.Chatbot) int64 {
	result := mdl.db.Create(&newChatbot)

	if result.Error != nil {
		log.Error(result.Error)
		return -1
	}

	return int64(newChatbot.ID)
}

func (mdl *model) SelectByID(chatbotID int) *chatbot.Chatbot {
	var chatbot chatbot.Chatbot
	result := mdl.db.First(&chatbot, chatbotID)

	if result.Error != nil {
		log.Error(result.Error)
		return nil
	}

	return &chatbot
}

func (mdl *model) Update(chatbot chatbot.Chatbot) int64 {
	result := mdl.db.Save(&chatbot)

	if result.Error != nil {
		log.Error(result.Error)
	}

	return result.RowsAffected
}

func (mdl *model) DeleteByID(chatbotID int) int64 {
	result := mdl.db.Delete(&chatbot.Chatbot{}, chatbotID)
	
	if result.Error != nil {
		log.Error(result.Error)
		return 0
	}

	return result.RowsAffected
}