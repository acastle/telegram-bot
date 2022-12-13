package command

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"telegram-bot/pkg/thread"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler interface {
	Exec(Context, *thread.Message) error
	IsReplyOnly() bool
}

type CommandRunner struct {
	handlers       map[string]Handler
	telegramClient *telegram.BotAPI
	gptClient      gpt3.Client
	repo           *thread.Repository
}

type Context struct {
	Runner   *CommandRunner
	Telegram *telegram.BotAPI
	GPT3     gpt3.Client
	Context  context.Context
	Threads  *thread.Repository
}

func (r *CommandRunner) SetTyping(chatID int64) error {
	action := telegram.NewChatAction(chatID, telegram.ChatTyping)
	_, err := r.telegramClient.Request(action)
	if err != nil {
		return fmt.Errorf("set typing status: %s", err)
	}

	return nil
}

func NewCommandRunnerFromEnv() (*CommandRunner, error) {
	telegramToken := os.Getenv("TELEGRAM_TOKEN")
	if telegramToken == "" {
		return nil, errors.New("no telegram token set")
	}

	openaiToken := os.Getenv("OPENAI_TOKEN")
	if openaiToken == "" {
		return nil, errors.New("no openai token set")
	}

	telegramClient, err := telegram.NewBotAPI(telegramToken)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	gptClient := gpt3.NewClient(openaiToken)
	var handlers = map[string]Handler{
		"echo":   Echo{},
		"prompt": Prompt{},
		"think":  Think{},
		"dump":   Dump{},
		"tweak":  Tweak{},
		"help":   Help{},
	}

	return &CommandRunner{
		handlers:       handlers,
		telegramClient: telegramClient,
		gptClient:      gptClient,
		repo:           thread.NewRepository(),
	}, nil
}

func (r *CommandRunner) Start() {
	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates := r.telegramClient.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() && update.Message.ReplyToMessage == nil {
			continue
		}

		r.DispatchHandler(update)
	}
}

func (r *CommandRunner) Reply(msg *thread.Message, text string, messageType thread.MessageType) error {
	reply := telegram.NewMessage(msg.ID.ChannelID, text)
	reply.ReplyToMessageID = msg.ID.MessageID
	newMsg, err := r.telegramClient.Send(reply)
	if err != nil {
		return fmt.Errorf("send prompt reply: %w", err)
	}

	_, err = r.repo.AddMessage(&newMsg, messageType)
	if err != nil {
		return fmt.Errorf("add message to thread: %w", err)
	}

	return nil
}

func (r *CommandRunner) DispatchHandler(update telegram.Update) {
	user := thread.User(*update.Message.From)
	log.Printf("Handle command [%s] %s", user.DisplayName(), update.Message.Text)
	var handler Handler = Prompt{}
	if update.Message.IsCommand() {
		cmd := update.Message.Command()
		var ok bool
		handler, ok = r.handlers[cmd]
		if !ok {
			log.Printf("skipping command '%s', no handler registered", cmd)
			return
		}
	}

	msgType := thread.TypeCommand
	if _, ok := handler.(Prompt); ok {
		msgType = thread.TypePrompt
	}

	msg, err := r.repo.AddMessage(update.Message, msgType)
	if err != nil {
		log.Printf("failed to add message to the repository")
	}

	if handler.IsReplyOnly() && update.Message.ReplyToMessage == nil {
		log.Printf("attempted to call a reply only handler without a reply message, skipping")
		r.Reply(msg, "That command can only be used in the context of a thread, try replying to an existing message.", thread.TypeInformational)
		return
	}

	timeout, _ := context.WithTimeout(context.Background(), 5*time.Minute)
	ctx := Context{
		Runner:   r,
		Telegram: r.telegramClient,
		GPT3:     r.gptClient,
		Context:  timeout,
		Threads:  r.repo,
	}

	ticker := time.NewTicker(3 * time.Second)
	done := make(chan struct{})
	go func(runner *CommandRunner, timeout <-chan struct{}, tick <-chan time.Time, done <-chan struct{}) {
		select {
		case <-timeout:
			ticker.Stop()
			log.Printf("command timed out")
			return
		case <-tick:
			if err := runner.SetTyping(update.Message.Chat.ID); err != nil {
				log.Printf("failed to set typing status: %s", err.Error())
			}
		case <-done:
			ticker.Stop()
			return
		}
	}(r, timeout.Done(), ticker.C, done)

	go func(runner *CommandRunner) {
		err := handler.Exec(ctx, msg)
		if err != nil {
			log.Printf("error in handler '%v': %s", handler, err.Error())
		}

		close(done)
	}(r)
}
