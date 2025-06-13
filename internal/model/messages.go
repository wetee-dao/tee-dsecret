package model

import (
	"errors"
	"fmt"
	"strings"

	"github.com/cockroachdb/pebble"
)

type BanTx struct {
	UserName string `json:"username"`
}

// Message represents a message sent by a user
type SideMessage struct {
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

type MsgHistory struct {
	Msg string `json:"history"`
}

func AppendToChat(message SideMessage) (string, error) {
	historyBytes, err := GetKey("", "history")
	if err != nil {
		fmt.Println("Error fetching history:", err)
		return "", err
	}
	msgBytes := string(historyBytes)
	msgBytes = msgBytes + "{sender:" + message.Sender + ",message:" + message.Message + "}"
	return msgBytes, nil
}

func FetchHistory() (string, error) {
	historyBytes, err := GetKey("", "history")
	if err != nil {
		fmt.Println("Error fetching history:", err)
		return "", err
	}
	msgHistory := string(historyBytes)

	if err != nil {
		fmt.Println("error appending history: ", err)
	}
	return msgHistory, err
}

func AppendToExistingMessages(message SideMessage) (string, error) {
	existingMessages, err := GetMessagesBySender(message.Sender)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		return "", err
	}
	if errors.Is(err, pebble.ErrNotFound) {
		return message.Message, nil
	}
	return existingMessages + ";" + message.Message, nil
}

// GetMessagesBySender retrieves all messages sent by a specific sender
// Get Message using String
func GetMessagesBySender(sender string) (string, error) {
	v, err := GetKey("", sender+"msg")
	if err != nil {
		return "", err
	}
	return string(v), nil
}

// ParseMessage parse messages
func ParseMessage(tx []byte) (*SideMessage, error) {
	msg := &SideMessage{}

	// Parse the message into key-value pairs
	pairs := strings.Split(string(tx), ",")

	if len(pairs) != 2 {
		return nil, errors.New("invalid number of key-value pairs in message")
	}

	for _, pair := range pairs {
		kv := strings.Split(pair, ":")

		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid key-value pair in message: %s", pair)
		}

		key := kv[0]
		value := kv[1]

		switch strings.ToLower(key) {
		case "sender":
			msg.Sender = value
		case "message":
			msg.Message = value
		case "history":
			return nil, fmt.Errorf("reserved key name: %s", key)
		default:
			return nil, fmt.Errorf("unknown key in message: %s", key)
		}
	}

	// Check if the message contains a sender and message
	if msg.Sender == "" {
		return nil, errors.New("message is missing sender")
	}

	if msg.Message == "" {
		return nil, errors.New("message is missing message")
	}

	return msg, nil
}
