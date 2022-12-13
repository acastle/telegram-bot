package thread

import (
	"reflect"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func TestUser_DisplayName(t *testing.T) {
	cases := []struct {
		desc     string
		input    User
		expected string
	}{
		{
			desc: "with first and last name",
			input: User{
				ID:        1234,
				FirstName: "Bob",
				LastName:  "TheBuilder",
			},
			expected: "Bob TheBuilder",
		},
		{
			desc: "with only first name",
			input: User{
				ID:        1234,
				FirstName: "Bob",
				LastName:  "",
			},
			expected: "Bob",
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			result := tc.input.DisplayName()
			if result != tc.expected {
				t.Errorf("expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestMessage_History(t *testing.T) {
	cases := []struct {
		desc     string
		input    *Message
		expected []string
	}{
		{
			desc: "includes prompts and responses in order",
			input: &Message{
				Text: "A response",
				Type: TypePrompt,
				Sender: User{
					FirstName: "Bob",
					LastName:  "TheBuilder",
				},
				parent: &Message{
					Text: "A prompt",
					Type: TypePrompt,
					Sender: User{
						FirstName: "Sam",
						LastName:  "TheClam",
					},
				},
			},
			expected: []string{
				"Sam TheClam: A prompt",
				"Bob TheBuilder: A response",
			},
		},
		{
			desc: "excludes informational messages and commands",
			input: &Message{
				Text: "Some command",
				Type: TypeCommand,
				Sender: User{
					FirstName: "Bob",
					LastName:  "TheBuilder",
				},
				parent: &Message{
					Text: "Some information",
					Type: TypeInformational,
					Sender: User{
						FirstName: "Suzie",
						LastName:  "TheSystem",
					},
					parent: &Message{
						Text: "A prompt",
						Type: TypePrompt,
						Sender: User{
							FirstName: "Sam",
							LastName:  "TheClam",
						},
					},
				},
			},
			expected: []string{
				"Sam TheClam: A prompt",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.desc, func(t *testing.T) {
			result := tc.input.History()
			if !reflect.DeepEqual(tc.expected, result) {
				t.Errorf("expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestGetMessageID_ProvidesDefault(t *testing.T) {
	result := GetMessageID(&tgbotapi.Message{
		MessageID: 1234,
		Chat:      nil,
		From:      nil,
	})
	expected := MessageID{MessageID: 1234}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected '%v', got '%v'", expected, result)
	}

	result = GetMessageID(&tgbotapi.Message{
		MessageID: 1234,
		Chat: &tgbotapi.Chat{
			ID: 5678,
		},
		From: nil,
	})
	expected = MessageID{MessageID: 1234, ChannelID: 5678}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected '%v', got '%v'", expected, result)
	}
}
