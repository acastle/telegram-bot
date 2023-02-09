package command

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"telegram-bot/pkg/thread"
)

type Tweak struct{}

const TweakParamHelp string = `/tweak [<parameter>=<value>;]
Parameters:
	Model:            <text-davinci-003|text-curie-001|text-babbage-001|text-ada-001>
	MaxTokens:        < 0 - 4000 >
	Temperature:      < 0.00 - 1.00 >
	FrequencyPenalty: < -2.00 - 2.00 >
	PressencePenalty: < -2.00 - 2.00 >
	TopP:             < 0.00 - 1.00 >
`

var ErrInvalidParameter error = errors.New("invalid parameter")

func (Tweak) Exec(ctx Context, msg *thread.Message) error {
	if msg.Parent() == nil {
		ctx.Runner.Reply(msg, "I couldn't find the thread history. I only keep a small sample of data around for future use, sorry 💩", thread.TypeInformational)
		return thread.ErrNotFound
	}

	currentThread, err := ctx.Threads.GetThread(msg.ThreadID)
	if err != nil {
		ctx.Runner.Reply(msg, "I couldn't find the thread history. I only keep a small sample of data around for future use, sorry 💩", thread.TypeInformational)
		return thread.ErrNotFound
	}

	var settings thread.CompletionParameters = currentThread.Settings
	setters := strings.Split(strings.TrimSuffix(msg.Text, ";"), ";")
	for _, setter := range setters {
		parts := strings.Split(strings.TrimSpace(setter), "=")
		if len(parts) != 2 {
			ctx.Runner.Reply(msg, fmt.Sprintf("Incorrect usage of command, the correct syntax is %s", TweakParamHelp), thread.TypeInformational)
			return ErrInvalidParameter
		}

		name := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		switch name {
		case "Model":
			settings, err = DoSet(settings, val, nil, setModel)
		case "MaxTokens":
			settings, err = DoSet(settings, val, strconv.Atoi, setMaxTokens)
		case "Temperature":
			settings, err = DoSet(settings, val, float32Converter, setTemperature)
		case "FrequencyPenalty":
			settings, err = DoSet(settings, val, float32Converter, setFrequencyPenalty)
		case "TopP":
			settings, err = DoSet(settings, val, float32Converter, setTopP)
		case "PressencePenalty":
			settings, err = DoSet(settings, val, float32Converter, setPressencePenalty)
		default:
			err = ErrInvalidParameter
		}

		if err != nil {
			ctx.Runner.Reply(msg, fmt.Sprintf("Incorrect usage of command, the correct syntax is %s", TweakParamHelp), thread.TypeInformational)
			return ErrInvalidParameter
		}

	}

	currentThread.Settings = settings
	ctx.Threads.Set(currentThread)
	return nil
}

type Params interface {
	float32 | int | string
}

func float32Converter(s string) (float32, error) {
	f, err := strconv.ParseFloat(s, 32)
	return float32(f), err
}

type converter[T Params] func(string) (T, error)
type paramSetter[T Params] func(thread.CompletionParameters, T) (thread.CompletionParameters, error)

func DoSet[T Params](params thread.CompletionParameters, valStr string, conv converter[T], setter paramSetter[T]) (thread.CompletionParameters, error) {
	var val T
	var err error
	if conv != nil {
		val, err = conv(valStr)
	}

	if err != nil {
		return params, err
	}

	return setter(params, val)
}

var ModelRegex = regexp.MustCompile("^(text-davinci-003|text-curie-001|text-babbage-001|text-ada-001)$")

func setModel(params thread.CompletionParameters, val string) (thread.CompletionParameters, error) {
	if !ModelRegex.Match([]byte(val)) {
		return params, ErrInvalidParameter
	}

	params.Model = val
	return params, nil
}

func setTopP(params thread.CompletionParameters, val float32) (thread.CompletionParameters, error) {
	if val < 0 || val > 1 {
		return params, ErrInvalidParameter
	}

	params.TopP = val
	return params, nil
}

func setPressencePenalty(params thread.CompletionParameters, val float32) (thread.CompletionParameters, error) {
	if val < -2 || val > 2 {
		return params, ErrInvalidParameter
	}

	params.PressencePenalty = val
	return params, nil
}

func setFrequencyPenalty(params thread.CompletionParameters, val float32) (thread.CompletionParameters, error) {
	if val < -2 || val > 2 {
		return params, ErrInvalidParameter
	}

	params.FrequencyPenalty = val
	return params, nil
}

func setTemperature(params thread.CompletionParameters, val float32) (thread.CompletionParameters, error) {
	if val < 0 || val > 1 {
		return params, ErrInvalidParameter
	}

	params.Temperature = val
	return params, nil
}

func setMaxTokens(params thread.CompletionParameters, val int) (thread.CompletionParameters, error) {
	if val < 0 || val > 4000 {
		return params, ErrInvalidParameter
	}

	params.MaxTokens = val
	return params, nil
}

func (Tweak) IsReplyOnly() bool {
	return true
}
