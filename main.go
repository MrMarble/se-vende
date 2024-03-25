package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Configure logger
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	if os.Getenv("DEBUG") == "true" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Logger = log.With().Caller().Logger()

	// Check required env vars
	if os.Getenv("TG_TOKEN") == "" {
		log.Fatal().Str("module", "main").Str("envvar", "TG_TOKEN").Msg("missing environment variable")
	}

	bot, err := NewBot(os.Getenv("TG_TOKEN"))
	if err != nil {
		log.Fatal().Str("module", "telegram").Err(err).Msg("failed bot instantiaion")
	}

	// Start bot
	bot.Start()
}
