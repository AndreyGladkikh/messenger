package chat

import (
	"fmt"
)

type ChatType string

const (
	PrivateChat = ChatType("private_chat")
	GroupChat = ChatType("group_chat")
) 

func CreateType(typ string) (ChatType, error) {
	switch typ {
	case "private_chat":
		return PrivateChat, nil
	case "group_chat":
		return GroupChat, nil
	default:
		return "", fmt.Errorf("unknown chat type")
	}
}