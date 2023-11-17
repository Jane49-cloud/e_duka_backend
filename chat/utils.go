package chat

import (
	"eleliafrika.com/backend/database"
)

func GetChats(conversation_id string) ([]Chat, error) {
	var chats []Chat
	err := database.Database.Where("conversation_id=?", conversation_id).Find(&chats).Error
	if err != nil {
		return []Chat{}, err
	}
	return chats, nil
}
