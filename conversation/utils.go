package conversation

import "eleliafrika.com/backend/database"

func FindClientConversation(client_id string) ([]Conversation, error) {
	var conversations []Conversation

	err := database.Database.Where("customer_id=?", client_id).Find(&conversations).Error
	if err != nil {
		return []Conversation{}, err
	}
	return conversations, nil
}
func FindSellerConversation(seller_id string) ([]Conversation, error) {
	var conversations []Conversation

	err := database.Database.Where("seller_id=?", seller_id).Find(&conversations).Error
	if err != nil {
		return []Conversation{}, err
	}
	return conversations, nil
}
