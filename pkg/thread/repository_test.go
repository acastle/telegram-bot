package thread

import (
	"errors"
	"reflect"
	"testing"

	"encoding/binary"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
)

var testingID uuid.UUID = uuid.UUID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xDE, 0xAD, 0xBE, 0xEF})

func generateTestingID() (uuid.UUID, error) {
	return testingID, nil
}

var offset uint32 = 0

func generateIncrementingTestID() (uuid.UUID, error) {
	b := [16]byte{}
	copy(b[:], testingID[:])

	val := binary.BigEndian.Uint32(b[12:])
	val += offset
	offset++

	binary.BigEndian.PutUint32(b[12:], val)
	return uuid.UUID(b), nil
}

func resetGenerator() {
	offset = 0
}

func TestRepository_NewThread(t *testing.T) {
	type output struct {
		value Thread
		err   error
	}
	cases := []struct {
		desc     string
		input    *Message
		expected output
		setup    func(*Repository)
	}{
		{
			desc: "fails on collisions",
			setup: func(repo *Repository) {
				generateID = generateTestingID
				repo.threads[testingID] = Thread{}
			},
			input: nil,
			expected: output{
				value: Thread{},
				err:   ErrAllocateThread,
			},
		},
		{
			desc: "retries on collision",
			setup: func(repo *Repository) {
				resetGenerator()
				generateID = generateIncrementingTestID
				repo.threads[testingID] = Thread{}
			},
			input: nil,
			expected: output{
				value: Thread{
					ID:       uuid.UUID([16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xDE, 0xAD, 0xBE, 0xF0}),
					Root:     nil,
					Settings: DefaultOpenAISettings,
				},
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			sut := NewRepository()
			tc.setup(sut)
			result, err := sut.NewThread(tc.input)
			if err == nil && !reflect.DeepEqual(result, tc.expected.value) {
				t.Errorf("expected '%v', got '%v'", tc.expected.value, result)
			}

			if !errors.Is(err, tc.expected.err) {
				t.Errorf("unexpected error expected '%v', got '%v'", tc.expected.err, err)
			}
		})
	}
}

func TestRepository_AddMessage(t *testing.T) {
	type input struct {
		message     telegram.Message
		messageType MessageType
	}
	type output struct {
		value *Message
		err   error
	}
	testUser := telegram.User{
		ID:        456,
		FirstName: "Foo",
		LastName:  "Bar",
	}
	cases := []struct {
		desc     string
		input    input
		expected output
		setup    func(*Repository)
	}{
		{
			desc: "creates new thread when one doesn't exist",
			setup: func(*Repository) {
				resetGenerator()
				generateID = generateTestingID
			},
			input: input{
				telegram.Message{
					MessageID: 1234,
					From:      &testUser,
					Text:      "Some message",
				},
				TypeInformational,
			},
			expected: output{
				value: &Message{
					ID: MessageID{
						ChannelID: 0,
						FromID:    456,
						MessageID: 1234,
					},
					ThreadID: testingID,
					Type:     TypeInformational,
					Text:     "Some message",
					Sender:   User(testUser),
					children: make([]*Message, 0),
					parent:   nil,
				},
				err: nil,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			sut := NewRepository()
			tc.setup(sut)
			result, err := sut.AddMessage(&tc.input.message, tc.input.messageType)
			if err == nil && result.ID != tc.expected.value.ID {
				t.Errorf("expected '%v', got '%v'", tc.expected.value, result)
			}

			if err == nil && result.Text != tc.expected.value.Text {
				t.Errorf("expected '%v', got '%v'", tc.expected.value, result)
			}

			if err == nil && result.Sender != tc.expected.value.Sender {
				t.Errorf("expected '%v', got '%v'", tc.expected.value, result)
			}

			if err == nil && result.Type != tc.expected.value.Type {
				t.Errorf("expected '%v', got '%v'", tc.expected.value, result)
			}

			if err == nil && result.ThreadID != tc.expected.value.ThreadID {
				t.Errorf("expected '%v', got '%v'", tc.expected.value, result)
			}

			if !errors.Is(err, tc.expected.err) {
				t.Errorf("unexpected error expected '%v', got '%v'", tc.expected.err, err)
			}
		})
	}
}
