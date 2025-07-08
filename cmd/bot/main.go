package main

import (
	"log"

	"github.com/vsespontanno/tgbot_fschedule/internal/bot"
)

func main() {
	err := bot.Start()
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
