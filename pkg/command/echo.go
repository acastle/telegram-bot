package command

import (
	"fmt"
	"telegram-bot/pkg/thread"
)

type Echo struct{}

func (Echo) Exec(ctx Context, msg *thread.Message) error {
	err := ctx.Runner.Reply(msg, msg.Text, thread.TypeInformational)
	if err != nil {
		return fmt.Errorf("send echo reply: %w", err)
	}

	return nil
}

func (Echo) IsReplyOnly() bool {
	return false
}
