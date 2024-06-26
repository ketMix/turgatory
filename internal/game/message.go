package game

import "image/color"

// MessageKind is the kind of message that is being sent.
type MessageKind int

const (
	// System Messages
	MessageDebug MessageKind = iota
	MessageInfo
	MessageError
	MessageUnknown

	// Game Messages
	MessageCombat
	MessageDamage
	MessageHeal
	MessageDeath
)

func (mk *MessageKind) Color() color.Color {
	switch *mk {
	case MessageDebug:
		return color.RGBA{150, 150, 150, 255}
	case MessageError:
		return color.RGBA{250, 75, 75, 255}
	case MessageCombat:
		return color.RGBA{75, 75, 250, 255}
	case MessageDamage:
		return color.RGBA{250, 75, 75, 255}
	case MessageHeal:
		return color.RGBA{75, 250, 75, 255}
	case MessageDeath:
		return color.RGBA{0, 0, 0, 255}
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
