package main

import (
	"football_tgbot/internal/bot"
	"log"
)

func main() {
	err := bot.Start()
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
}
