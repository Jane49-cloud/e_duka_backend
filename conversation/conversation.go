package conversation

import "eleliafrika.com/backend/database"

type Conversation struct {
	ConversationId string `gorm:"column:conversation_id;no null;unique" json:"conversation_id"`
	CustomerId     string `gorm:"customer_id;not null;require" json:"customer_id"`
	SellerId       string `gorm:"seller_id;not null;require" json:"seller_id"`
	LastMessage    string `gorm:"column:last_message;type:text;" json:"last_message"`
	MessagesCount  int    `gorm:"column:messages_count;default:0;" json:"messages_count"`
	IsViewed       bool   `gorm:"column:is_viewed;default:false;" json:"is_viewed"`
	IsDeleted      bool   `gorm:"column:is_deleted;default:false;" json:"is_deleted"`
	DateAdded      string `gorm:"column:date_created;not null" json:"date_added"`
}

type ConversationInput struct {
	Message  string `gorm:"message;not null;require" json:"message"`
	SellerId string `gorm:"seller_id;not null;require" json:"seller_id"`
}

func (conversation *Conversation) Save() (*Conversation, error) {
	err := database.Database.Create(&conversation).Error
	if err != nil {
		return &Conversation{}, err
	}
	return conversation, nil
}
