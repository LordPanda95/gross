package internal

import (
	"fmt"

	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/rs/zerolog"
)

func StartHandler(bh *th.BotHandler, chatIdInt64 int64, config *AppConfig, logger *zerolog.Logger) error {
	handlerLogger := logger.With().Str("func", "StartHandler").Logger()
	// Register new handler with match on command `/gross` and args `today`
	bh.Handle(func(bot *telego.Bot, update telego.Update) {

		// Get current time
		currentTime := time.Now()
		handlerLogger.Debug().
			Msgf("currentTime", currentTime.GoString())
		today := currentTime.Format("2006-01-02")
		handlerLogger.Debug().
			Msg("today: " + today)

		// Get gross massage
		grossMessage, err := MessageText(today, config.GrossConfig, config.Nbrb, logger)
		if err != nil {
			handlerLogger.Error().
				Err(err).
				Str("grossMessage", grossMessage).
				Msg("cannot getting gross message")
			return
		}

		// Send message
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(chatIdInt64),
			fmt.Sprintf("Hello %s! ITS GROSS TODAY"+grossMessage, update.Message.From.FirstName),
		))
	}, th.CommandEqual("gross"), th.TextContains("today"))

	// Register new handler with match on command `/gross` and args `tomorrow`
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Get current time
		currentTime := time.Now()
		handlerLogger.Debug().
			Msgf("currentTime", currentTime.GoString())
		tomorrow := currentTime.AddDate(0, 0, 1).Format("2006-01-02")
		handlerLogger.Debug().
			Msg("tomorrow: " + tomorrow)

		// Get gross massage
		grossMessage, err := MessageText(tomorrow, config.GrossConfig, config.Nbrb, logger)
		if err != nil {
			handlerLogger.Error().
				Err(err).
				Str("grossMessage", grossMessage).
				Msg("cannot getting gross message")
			return
		}
		// Send message
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(chatIdInt64),
			fmt.Sprintf("Hello %s! ITS GROSS TOMORROW"+grossMessage, update.Message.From.FirstName),
		))
	}, th.CommandEqual("gross"), th.TextContains("tomorrow"))

	// Register new handler with match on command `/gross`
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Get current time
		currentTime := time.Now()
		handlerLogger.Debug().
			Msgf("currentTime", currentTime.GoString())
		today := currentTime.Format("2006-01-02")
		handlerLogger.Debug().
			Msg("today: " + today)

		// Get gross massage
		grossMessage, err := MessageText(today, config.GrossConfig, config.Nbrb, logger)
		if err != nil {
			handlerLogger.Error().
				Err(err).
				Str("grossMessage", grossMessage).
				Msg("cannot getting gross message")
			return
		}

		// Send message
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(chatIdInt64),
			fmt.Sprintf("Hello %s!"+grossMessage, update.Message.From.FirstName),
		))
	}, th.CommandEqual("gross"))

	// Register new handler with match on command `/start`
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Send message
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(chatIdInt64),
			fmt.Sprintf("Hello %s!", update.Message.From.FirstName),
		))
	}, th.CommandEqual("start"))

	// Register new handler with match on any command
	// Handlers will match only once and in order of registration,
	// so this handler will be called on any command except `/start` command
	bh.Handle(func(bot *telego.Bot, update telego.Update) {
		// Send message
		_, _ = bot.SendMessage(tu.Message(
			tu.ID(chatIdInt64),
			"Unknown command, use /start",
		))
	}, th.AnyCommand())

	return nil
}
