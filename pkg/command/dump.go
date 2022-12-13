package command

import (
	"fmt"
	"strings"
	"telegram-bot/pkg/thread"
)

type Dump struct{}

func (Dump) Exec(ctx Context, msg *thread.Message) error {
	if msg.Parent() == nil {
		ctx.Runner.Reply(msg, "I couldn't find the thread history. I only keep a small sample of data around for future use, sorry 💩", thread.TypeInformational)
		return thread.ErrNotFound
	}

	t, err := ctx.Threads.GetThread(msg.ThreadID)
	if err != nil {
		ctx.Runner.Reply(msg, "I couldn't find the thread history. I only keep a small sample of data around for future use, sorry 💩", thread.TypeInformational)
		return thread.ErrNotFound
	}

	ctx.Runner.Reply(msg, BuildReportForThread(t), thread.TypeInformational)
	return nil
}

func BuildReportForThread(t thread.Thread) string {
	var b strings.Builder
	b.WriteString("Thread Stats:\n")
	b.WriteString(fmt.Sprintf("  Thread ID: %s\n", t.ID))
	b.WriteString(fmt.Sprintf("    Total messages:\t\t%d\n", SumChildren(*t.Root, IncludeAllMessages)))
	b.WriteString(fmt.Sprintf("    Total prompts:\t\t%d\n", SumChildren(*t.Root, IncludeChildrenOfType(thread.TypePrompt))))
	b.WriteString(fmt.Sprintf("    Total responses:\t\t%d\n", SumChildren(*t.Root, IncludeChildrenOfType(thread.TypeResponse))))
	b.WriteString(fmt.Sprintf("    Total commands:\t\t%d\n", SumChildren(*t.Root, IncludeChildrenOfType(thread.TypeCommand))))
	b.WriteString(fmt.Sprintf("    Total informational:\t\t%d\n", SumChildren(*t.Root, IncludeChildrenOfType(thread.TypeInformational))))
	b.WriteString("\n")
	b.WriteString("Thread GPT3 Parameters:\n")
	b.WriteString(fmt.Sprintf("    MaxTokens:\t\t%d\n", t.Settings.MaxTokens))
	b.WriteString(fmt.Sprintf("    FrequencyPenalty:\t\t%f\n", t.Settings.FrequencyPenalty))
	b.WriteString(fmt.Sprintf("    PressencePenalty:\t\t%f\n", t.Settings.PressencePenalty))
	b.WriteString(fmt.Sprintf("    Temperature:\t\t%f\n", t.Settings.Temperature))
	b.WriteString(fmt.Sprintf("    TopP:\t\t%f\n", t.Settings.TopP))
	return b.String()
}

type Predicate = func(thread.Message) bool

func IncludeAllMessages(msg thread.Message) bool {
	return true
}

func IncludeChildrenOfType(t thread.MessageType) Predicate {
	return func(msg thread.Message) bool {
		return msg.Type == t
	}
}

func SumChildren(msg thread.Message, predicate Predicate) int {
	childTotal := 0
	for _, c := range msg.Children() {
		childTotal += SumChildren(*c, predicate)
	}

	if predicate(msg) {
		return childTotal + 1
	}

	return childTotal
}

func (Dump) IsReplyOnly() bool {
	return true
}
