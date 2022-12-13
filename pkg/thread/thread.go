package thread

import (
	"errors"
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

var (
	ErrAllocateThread = errors.New("failed to allocate thread")
	ErrNotFound       = errors.New("could not find element")
)

type Thread struct {
	ID       uuid.UUID
	Root     *Message
	Settings CompletionParameters
}

type MessageID struct {
	ChannelID int64
	FromID    int64
	MessageID int
}

type MessageType int

const (
	TypePrompt MessageType = iota
	TypeResponse
	TypeInformational
	TypeCommand
)

type Message struct {
	ID       MessageID
	Type     MessageType
	ThreadID uuid.UUID
	Sender   User
	Text     string

	parent   *Message
	children []*Message
}

type User telegram.User

func (u *User) DisplayName() string {
	if u.LastName == "" {
		return u.FirstName
	}

	return fmt.Sprintf("%s %s", u.FirstName, u.LastName)
}

func GetMessageID(msg *telegram.Message) MessageID {
	var channelID int64
	var fromID int64

	if msg.Chat != nil {
		channelID = msg.Chat.ID
	}

	if msg.From != nil {
		fromID = msg.From.ID
	}

	return MessageID{
		ChannelID: channelID,
		FromID:    fromID,
		MessageID: msg.MessageID,
	}
}

func (m *Message) History() []string {
	prompt := fmt.Sprintf("%s: %s", m.Sender.DisplayName(), m.Text)
	history := []string{}
	if m.parent != nil {
		history = m.parent.History()
	}

	if m.Type == TypePrompt || m.Type == TypeResponse {
		return append(history, prompt)
	}

	return history
}

func (m *Message) Parent() *Message {
	return m.parent
}

func (m *Message) Children() []*Message {
	return m.children
}
