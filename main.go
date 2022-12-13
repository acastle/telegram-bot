package main

import (
	"log"
	"telegram-bot/pkg/command"
)

var handlers = map[string]command.Handler{
	"echo": command.Echo{},
}

func main() {
	log.Println("Starting telegram bot")
	runner, err := command.NewCommandRunnerFromEnv()
	if err != nil {
		panic(err)
	}

	runner.Start()
}
