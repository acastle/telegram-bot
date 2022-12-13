package command

import (
	"fmt"
	"strings"
	"telegram-bot/pkg/thread"

	"github.com/PullRequestInc/go-gpt3"
)

type Prompt struct{}

func (Prompt) Exec(ctx Context, msg *thread.Message) error {
	stopTokens := []string{}
	currentThread, err := ctx.Threads.GetThread(msg.ThreadID)
	if err != nil {
		return fmt.Errorf("get thread: %w", err)
	}

	prompts := msg.History()
	bot := thread.User(ctx.Telegram.Self)
	prompts = append(prompts, bot.DisplayName()+":")
	prompt := strings.Join(prompts, "\n\n")
	resp, err := ctx.GPT3.CompletionWithEngine(ctx.Context, currentThread.Settings.Model, gpt3.CompletionRequest{
		Prompt:           []string{prompt},
		MaxTokens:        &currentThread.Settings.MaxTokens,
		Temperature:      &currentThread.Settings.Temperature,
		FrequencyPenalty: currentThread.Settings.FrequencyPenalty,
		PresencePenalty:  currentThread.Settings.PressencePenalty,
		TopP:             &currentThread.Settings.TopP,
		Stop:             stopTokens,
	})
	if err != nil {
		return fmt.Errorf("call gpt3 api: %w", err)
	}

	ctx.Runner.Reply(msg, resp.Choices[0].Text, thread.TypeResponse)
	return nil
}

func (Prompt) IsReplyOnly() bool {
	return false
}
