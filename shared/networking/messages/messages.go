package messages

import (
	"go.mongodb.org/mongo-driver/bson"
)

type Message struct {
	Content string      `bson:"content"`
	Type    messageType `bson:"type"`
}

type messageType string

const (
	Unknown    messageType = "unknown"
	Normal     messageType = "normal"
	Connect    messageType = "connect"
	Disconnect messageType = "disconnect"
	Closed     messageType = "closed"
	Heartbeat  messageType = "heartbeat"
)

func New(content string, t messageType) *Message {
	return &Message{
		Content: content,
		Type:    t,
	}
}

func (m *Message) Serialize() []byte {
	// convert to bson
	data, err := bson.Marshal(m)
	if err != nil {
		return nil
	}

	return data
}

func Deserialize(data []byte) *Message {
	var m Message
	err := bson.Unmarshal(data, &m)
	if err != nil {
		return nil
	}

	return &m
}
