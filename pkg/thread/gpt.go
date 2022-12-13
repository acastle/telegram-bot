package thread

type CompletionParameters struct {
	Model            string
	MaxTokens        int
	Temperature      float32
	FrequencyPenalty float32
	PressencePenalty float32
	TopP             float32
}

var DefaultOpenAISettings = CompletionParameters{
	Model:            "text-davinci-003",
	MaxTokens:        400,
	Temperature:      0.5,
	FrequencyPenalty: 0,
	PressencePenalty: 0,
	TopP:             1,
}
