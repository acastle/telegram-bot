package command

import (
	"fmt"
	"telegram-bot/pkg/thread"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Think struct{}

func (Think) Exec(ctx Context, msg *thread.Message) error {
	time.Sleep(5 * time.Second)
	reply := telegram.NewMessage(msg.ID.ChannelID, "had a good sleep")
	if _, err := ctx.Telegram.Send(reply); err != nil {
		return fmt.Errorf("send think reply: %w", err)
	}

	return nil
}

func (Think) IsReplyOnly() bool {
	return false
}
