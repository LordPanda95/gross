package shared

import (
	"errors"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	"github.com/rs/zerolog"
)

type Telegram struct {
	Bot     *telego.Bot
	Handler *telegohandler.BotHandler
	Config  *TelegramConfig
}

type TelegramConfig struct {
	ApiKey string `mapstructure:"api_key"`
	ChatId int64  `mapstructure:"chat_id"`
}

func (telegram *Telegram) NewBot(config *TelegramConfig, loggerConfig *LoggerConfig, logger *zerolog.Logger) error {
	// Note: Please keep in mind that default logger may expose sensitive information,

	err := errors.New("")

	// use in development only
	if loggerConfig.Level == "debug" {
		telegram.Bot, err = telego.NewBot(config.ApiKey, telego.WithDefaultDebugLogger())
		logger.Info().
			Msg("telegram bot created in debug mode")
		if err != nil {
			logger.Fatal().
				Err(err).
				Msg("cannot create bot")
		}

	} else {
		telegram.Bot, err = telego.NewBot(config.ApiKey)
		logger.Info().
			Msg("telegram bot created")
		if err != nil {
			logger.Fatal().
				Err(err).
				Msg("cannot create bot")
		}
	}

	telegram.Bot, err = telego.NewBot(config.ApiKey)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("cannot create bot")
	}

	// Get updates channel
	updates, _ := telegram.Bot.UpdatesViaLongPolling(nil)

	// Create bot handler and specify from where to get updates
	telegram.Handler, err = telegohandler.NewBotHandler(telegram.Bot, updates)
	if err != nil {
		logger.Fatal().
			Err(err).
			Msg("cannot create bot handler")
	}

	return nil
}
