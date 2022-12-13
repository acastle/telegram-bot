package thread

import (
	"fmt"
	"sync"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

const IDCollisionRetryCount int = 5

type Repository struct {
	threads  map[uuid.UUID]Thread
	messages map[MessageID]*Message

	mu sync.Mutex
}

func NewRepository() *Repository {
	return &Repository{
		threads:  make(map[uuid.UUID]Thread),
		messages: map[MessageID]*Message{},
	}
}

func (r *Repository) AddMessage(source *telegram.Message, messageType MessageType) (*Message, error) {
	var parent *Message
	id := GetMessageID(source)
	msg := Message{
		ID:       id,
		Type:     messageType,
		parent:   nil,
		children: []*Message{},
		Text:     source.Text,
	}

	if source.From != nil {
		msg.Sender = User(*source.From)
	}

	if source.IsCommand() {
		msg.Text = source.CommandArguments()
	}

	if source.ReplyToMessage != nil {
		id := GetMessageID(source.ReplyToMessage)
		parent = r.GetMessage(id)
		msg.parent = parent
	}

	var threadID uuid.UUID
	if parent == nil {
		thread, err := r.NewThread(&msg)
		if err != nil {
			return nil, fmt.Errorf("allocate thread: %w", err)
		}

		threadID = thread.ID
	} else {
		parent.children = append(parent.children, &msg)
		threadID = parent.ThreadID
	}

	msg.ThreadID = threadID
	r.mu.Lock()
	r.messages[id] = &msg
	r.mu.Unlock()
	return &msg, nil
}

func (r *Repository) Set(thread Thread) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.threads[thread.ID] = thread
}

var generateID func() (uuid.UUID, error) = uuid.NewUUID

func (r *Repository) NewThread(root *Message) (t Thread, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	var id uuid.UUID
	for i := 0; i < IDCollisionRetryCount; i++ {
		id, err = generateID()
		if err != nil {
			continue
		}

		_, ok := r.threads[id]
		if !ok {
			t = Thread{
				ID:       id,
				Root:     root,
				Settings: DefaultOpenAISettings,
			}

			r.threads[id] = t
			return t, nil
		}
	}

	return t, ErrAllocateThread
}

func (r *Repository) GetMessage(id MessageID) *Message {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.messages[id]
}

func (r *Repository) GetThread(id uuid.UUID) (t Thread, err error) {
	r.mu.Lock()
	var ok bool
	t, ok = r.threads[id]
	r.mu.Unlock()
	if ok {
		return t, nil
	}

	return t, ErrNotFound
}
