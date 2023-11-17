package chat

import (
	"eleliafrika.com/backend/database"
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	ChatID         string `gorm:"column:chat_id;not null;primary key;unique;" json:"chat_id"`
	ConversationId string `gorm:"column:conversation_id;not null;primary key;" json:"conversation_id"`
	SenderID       string `gorm:"column:sender_id;size:255;not null;" json:"sender_id"`
	ReceiverId     string `gorm:"column:seller_id;size:255;not null;" json:"receiver_id"`
	Message        string `gorm:"column:message_body;type:text;not null;" json:"message_body"`
	TimeSent       string `gorm:"column:time_sent;not null;" json:"time_sent"`
	IsViewed       bool   `gorm:"column:is_viewed;default:false" json:"is_viewed"`
}
type ChatInput struct {
	ReceiverId string `json:"receiver_id"`
	Message    string `json:"message_body"`
}

func (chat *Chat) Save() (*Chat, error) {
	err := database.Database.Create(&chat).Error
	if err != nil {
		return &Chat{}, err
	}
	return chat, nil
}
