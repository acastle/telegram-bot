# telegram-bot
WIP GPT3 bot for telegram

## Running

Run with `go run main.go`

## Configuration

All configuration provided with environment variables

| Option | Description | Required |
|---|---|---|
| `OPENAI_TOKEN` | Access token for the OpenAI API. Register at https://beta.openai.com/account/api-keys | x |
| `TELEGRAM_TOKEN` | Bot token for Telegram. Generated via messaging @botfather `/newbot` | x |

## Supported commands


| Command | Description | Reply Only |
|---|---|---|
| `/prompt <text>` | Starts a new thread using the provided text as a seed. |  |
| `/echo <text>` | Echos the provided text as a response. Useful for being a new thread without a prompt. |  |
| `/think` | Doesn't do anything particularly useful, waits for 5 seconds and replies. To be removed in the future. |  |
| `/dump` | Responds with information about the current chat thread, number of prompts, responses, informational messages, etc. As well as the current thread's GPT parameters. | x |
| `/tweak [<param>=<value>;]` | <table><br><thead><br><tr><br><th>Parameter</th><br><th>Value</th><br></tr><br></thead><br><tbody><br><tr><br><td>Model</td><br><td><code>text-davinci-003\|text-curie-001\|text-babbage-001\|text-ada-001</code></td><br></tr><br><tr><br><td>MaxTokens</td><br><td><code>0 - 4000</code></td><br></tr><br><tr><br><td>Temperature</td><br><td><code>0.00 - 1.00</code></td><br></tr><br><tr><br><td>FrequencyPenalty</td><br><td><code>-2.00 - 2.00</code></td><br></tr><br><tr><br><td>PressencePenalty</td><br><td><code>-2.00 - 2.00</code></td><br></tr><br><tr><br><td>TopP</td><br><td><code>0.00 - 1.00</code></td><br></tr><br></tbody><br></table> | x |
