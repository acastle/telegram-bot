package command

import (
	"strings"
	"telegram-bot/pkg/thread"
)

type Help struct{}

func (Help) Exec(ctx Context, msg *thread.Message) error {
	ctx.Runner.Reply(msg, PrintHelp(), thread.TypeInformational)
	return nil
}

func PrintHelp() string {
	var b strings.Builder
	b.WriteString("```\n")
	b.WriteString("Commands:\n")
	b.WriteString("  /prompt <text>: Initate a new thread starting with the provided prompt.\n")
	b.WriteString("  /echo <text>:   Reply with the exact text (starts a new thread without prompt)\n")
	b.WriteString("  /help:          Prints this text\n")
	b.WriteString("\n")
	b.WriteString("Reply only commands:\n")
	b.WriteString("  /dump: Dumps out technical information about the current conversation thread\n")
	b.WriteString(TweakParamHelp)
	b.WriteString("```\n")
	return b.String()
}

func (Help) IsReplyOnly() bool {
	return false
}
