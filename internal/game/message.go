package game

import "image/color"

// MessageKind is the kind of message that is being sent.
type MessageKind int

const (
	// System Messages
	MessageDebug MessageKind = iota
	MessageUnknown

	MessageNeutral
	MessageLoot
	MessageGood
	MessageBad
)

func (mk *MessageKind) Color() color.Color {
	switch *mk {
	case MessageNeutral:
		return color.RGBA{125, 125, 125, 255}
	case MessageGood:
		return color.RGBA{75, 250, 75, 255}
	case MessageBad:
		return color.RGBA{250, 75, 75, 255}
	case MessageLoot:
		return color.RGBA{255, 215, 0, 255}
	default:
		return color.RGBA{200, 200, 200, 255}
	}

}

type Message struct {
	kind MessageKind
	text string
}

var messages []*Message

// AddMessage adds a message to the message log
func AddMessage(kind MessageKind, text string) *Message {
	m := &Message{
		kind: kind,
		text: text,
	}
	messages = append(messages, m)

	// Limit the number of messages to 500
	if len(messages) > 500 {
		messages = messages[1:]
	}
	return m
}

func GetMessages() []*Message {
	return messages
}
