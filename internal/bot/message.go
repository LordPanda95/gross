package internal

import (
	"fmt"

	"gross/pkg/shared"

	"github.com/rs/zerolog"
)

func getMessageData(date string, grossConfig *GrossConfig, nbrb *shared.Nbrb, logger *zerolog.Logger) (*shared.Currency, *Gross, error) {

	grossLogger := logger.With().Str("func", "getGross").Logger()

	grossLogger.Debug().
		Msgf("getting gross for date: %v", date)

	// Get currency rates
	currencyRate, err := nbrb.GetCurrencyRate(date)
	if err != nil {
		grossLogger.Error().
			Err(err).
			Str("currencyRate", fmt.Sprintf("%v", currencyRate)).
			Msg("cannot getting currency rates")
		return &shared.Currency{}, &Gross{}, err
	}

	// Calculate the gross summary
	gross := &Gross{}
	gross.Get(currencyRate, grossConfig, logger)

	grossLogger.Debug().
		Str("gross summary", fmt.Sprintf("%v", gross)).
		Str("currencyRate", fmt.Sprintf("%v", currencyRate)).
		Msg("gross summary calculated")

	return currencyRate, gross, nil
}
func MessageText(date string, grossConfig *GrossConfig, nbrb *shared.Nbrb, logger *zerolog.Logger) (string, error) {

	messageLogger := logger.With().Str("func", "GrossMessageText").Logger()

	messageLogger.Debug().
		Msgf("getting gross message for date: %v", date)

	currency, gross, err := getMessageData(date, grossConfig, nbrb, logger)
	if err != nil {
		messageLogger.Error().
			Err(err).
			Str("grossData", fmt.Sprintf("%v", gross)).
			Msg("cannot getting gross message")
		return "", err
	}

	// Build message
	grossMessage := fmt.Sprintln("")
	grossMessage = grossMessage + fmt.Sprintln("")
	grossMessage = grossMessage + fmt.Sprintln("Date:", date)
	grossMessage = grossMessage + fmt.Sprintln("")
	grossMessage = grossMessage + fmt.Sprintln("Курсы валют:")
	grossMessage = grossMessage + fmt.Sprintln("")
	grossMessage = grossMessage + fmt.Sprintln(currency.Rub, "Белорусских рублей за", currency.RubScale, "Российских рублей")
	grossMessage = grossMessage + fmt.Sprintln(currency.Usd, "Белорусских рублей за 1 Доллар США")
	grossMessage = grossMessage + fmt.Sprintln("")
	grossMessage = grossMessage + fmt.Sprintln("Зарплата в рублях:", gross.Fact.Gross)
	grossMessage = grossMessage + fmt.Sprintln("Зарплата в рублях после уплаты налога:", gross.Fact.Net)
	grossMessage = grossMessage + fmt.Sprintln("Налоги в рублях:", gross.Fact.Tax)
	grossMessage = grossMessage + fmt.Sprintln("")
	grossMessage = grossMessage + fmt.Sprintln("Планируемая зарплата в рублях:", gross.Plan.Gross)
	grossMessage = grossMessage + fmt.Sprintln("Планируемая зарплата в рублях после уплаты налога:", gross.Plan.Net)
	grossMessage = grossMessage + fmt.Sprintln("Планируемые налоги в рублях:", gross.Plan.Tax)
	grossMessage = grossMessage + fmt.Sprintln("")
	grossMessage = grossMessage + fmt.Sprintln("Итоговый профит:", gross.Profit)

	messageLogger.Debug().
		Str("grossMessage", grossMessage).
		Msg("gross message")

	return grossMessage, nil
}
