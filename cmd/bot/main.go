package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	internal "gross/internal/bot"
	"gross/pkg/shared"

	"github.com/rs/zerolog/log"
)

func main() {

	// Load the config from the .env file to the config struct.
	config := internal.AppConfig{}
	err := shared.LoadConfig(".", &config)
	log.Info().
		Str("func", "main").
		Msg("loading config...")
	if err != nil {
		log.Fatal().
			Str("func", "main").
			Err(err).
			Msg("error loading config")
	} else {
		log.Info().
			Str("func", "main").
			Msg("config loaded successfully")
	}

	logger, err := shared.NewLogger(config.Logger)
	if err != nil {
		log.Fatal().
			Str("func", "main").
			Err(err).
			Msg("error initializing logger")
	}

	mainLogger := logger.With().Str("package", "main").Logger()

	// Check if the config struct is empty.
	err = shared.ConfigIsEmpty(config, "config", logger)
	mainLogger.Debug().
		Str("config", fmt.Sprintf("%#v", config)).
		Msg("checking if config is empty...")
	if err != nil {
		mainLogger.Error().
			Err(err).
			Msg("")
	}
	mainLogger.Debug().
		Msg("checking if config is empty successful")

	// Create a channel to receive the OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		telegram := shared.Telegram{}
		err = telegram.NewBot(config.Telegram, config.Logger, logger)
		if err != nil {
			mainLogger.Error().
				Err(err).
				Msg("cannot create bot handler")
		}

		// Stop handling updates
		defer telegram.Handler.Stop()

		// Stop getting updates
		defer telegram.Bot.StopLongPolling()

		// Start getting updates
		err = internal.StartHandler(telegram.Handler, config.Telegram.ChatId, &config, logger)
		if err != nil {
			mainLogger.Error().
				Err(err).
				Msg("cannot start handler")
		}

		// Start handling updates
		telegram.Handler.Start()
	}()

	// Start a goroutine to handle the SIGTERM signal
	//go func() {
	//	<-sigCh
	//	fmt.Println("Received SIGTERM signal. Starting graceful shutdown...")
	//	// Perform any necessary cleanup operations here
	//
	//	// Cancel the context to signal the goroutines to stop
	//	cancel()
	//}()

	sig := <-sigCh
	fmt.Println("Received SIGTERM signal. Starting graceful shutdown...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	<-ctx.Done()

	fmt.Println("Graceful shutdown completed.")
}
