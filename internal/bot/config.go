package internal

import "gross/pkg/shared"

type AppConfig struct {
	Telegram    *shared.TelegramConfig `mapstructure:"telegram"`
	Nbrb        *shared.Nbrb           `mapstructure:"nbrb"`
	GrossConfig *GrossConfig           `mapstructure:"gross"`
	Logger      *shared.LoggerConfig   `mapstructure:"logger"`
}
