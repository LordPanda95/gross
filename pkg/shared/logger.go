package shared

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// Create struct for config slack.
type LoggerConfig struct {
	Service string `mapstructure:"service"`
	Level   string `mapstructure:"level"`
}

func NewLogger(config *LoggerConfig) (*zerolog.Logger, error) {

	// Set the logging level for the zerolog package
	switch config.Level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Enable error stack marshaling
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	// Create a new logger with default fields
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "vacation_bot").
		//Caller().
		Stack().
		Logger()

	return &logger, nil
}
