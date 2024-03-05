package internal

import (
	"errors"
	"math"
	"strconv"

	"gross/pkg/shared"

	"github.com/rs/zerolog"
)

type Gross struct {
	Fact   *FactGross
	Plan   *PlanGross
	Profit float64
}

type FactGross struct {
	Gross float64
	Net   float64
	Tax   float64
}

type PlanGross struct {
	Gross float64
	Net   float64
	Tax   float64
}

type GrossConfig struct {
	GrossUsd    float64 `mapstructure:"gross_usd"`
	TaxRate     float64 `mapstructure:"tax_rate"`
	PlanUsdRate float64 `mapstructure:"plan_usd_rate"`
}

func (gross *Gross) Get(currency *shared.Currency, grossConfig *GrossConfig, logger *zerolog.Logger) error {

	err := errors.New("")

	gross.Fact, err = getFactGross(currency, grossConfig, logger)
	if err != nil {
		return err
	}

	gross.Plan, err = getPlanGross(grossConfig, logger)
	if err != nil {
		return err
	}

	// Calculate the profit
	gross.Profit, err = getProfit(gross.Fact.Net, gross.Plan.Net)
	if err != nil {
		return err
	}

	// Return the gross
	return nil
}

func getProfit(realNet float64, planNet float64) (float64, error) {

	// Calculate the profit
	profit := realNet - planNet

	// До 2 знаков после запятой
	profit = math.Round(profit*100) / 100

	return profit, nil
}

func getFactGross(currency *shared.Currency, grossConfig *GrossConfig, logger *zerolog.Logger) (*FactGross, error) {

	calculateLogger := logger.With().Str("func", "Calculate").Logger()

	calculateLogger.Debug().
		Str("rateUsd", strconv.FormatFloat(currency.Usd, 'f', -1, 64)).
		Str("rateRub", strconv.FormatFloat(currency.Rub, 'f', -1, 64)).
		Str("rateRubScale", strconv.FormatFloat(currency.RubScale, 'f', -1, 64)).
		Msgf("calculating gross summary")

	// Зарплата грязными
	gross := ((grossConfig.GrossUsd * currency.Usd) / currency.Rub) * currency.RubScale

	// Округление до 2 знаков после запятой
	gross = math.Round(gross*100) / 100

	// Налоги
	tax := (gross / 100) * grossConfig.TaxRate

	// Округление до 2 знаков после запятой
	tax = math.Round(tax*100) / 100

	// Зарплата чистыми
	net := gross - tax

	// Округление до 2 знаков после запятой
	net = math.Round(net*100) / 100

	// Return the gross summary
	return &FactGross{
		Gross: gross,
		Net:   net,
		Tax:   tax,
	}, nil

}

func getPlanGross(grossConfig *GrossConfig, logger *zerolog.Logger) (*PlanGross, error) {

	// !!
	//calculateLogger := logger.With().Str("func", "Calculate").Logger()

	// Планируемая зарплата
	gross := grossConfig.GrossUsd * grossConfig.PlanUsdRate

	// Округление до 2 знаков после запятой
	gross = math.Round(gross*100) / 100

	// Плановые Налоги
	tax := (gross / 100) * grossConfig.TaxRate

	// Округление до 2 знаков после запятой
	tax = math.Round(tax*100) / 100

	// Зарплата чистыми
	net := gross - tax

	// Округление до 2 знаков после запятой
	net = math.Round(net*100) / 100

	// Return the gross summary
	return &PlanGross{
		Gross: gross,
		Net:   net,
		Tax:   tax,
	}, nil
}
