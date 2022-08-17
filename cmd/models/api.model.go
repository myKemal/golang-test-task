package models

type NewMessagePayload struct {
	Sender   string `json:"sender"`
	Receiver string `json:"reciever"`
	Message  string `json:"message"`
}

type GetMessagePayload struct {
	Sender   string `json:"sender"`
	Receiver string `json:"reciever"`
}

type MessageList struct {
	Message []string `json:"message"`
}

type RabbitMessage struct {
	Quene   string            `json:"quene"`
	Message NewMessagePayload `json:"message"`
}
