package shared

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func ConfigIsEmpty(structure interface{}, path string, logger *zerolog.Logger) error {
	// Add new default fields to the logger for this function.
	configEmptyLogger := logger.With().Str("func", "ConfigIsEmpty").Logger()

	configEmptyLogger.Debug().
		Str("path", path).
		Str("struct", fmt.Sprintf("%#v", structure)).
		Msg("checking if config is empty...")

	// Get the value of the structure.
	structValue := reflect.ValueOf(structure)
	configEmptyLogger.Debug().
		Str("structValue", fmt.Sprintf("%#v", structValue)).
		Msg("getting struct value...")
	// Get the type of the structure.
	structType := structValue.Type()
	configEmptyLogger.Debug().
		Str("structType", fmt.Sprintf("%#v", structType)).
		Msg("getting struct type...")

	// Check if the structure is empty.
	for i := 0; i < structValue.NumField(); i++ {
		// Get the field in the structure.
		field := structValue.Field(i)
		configEmptyLogger.Debug().
			Str("field", fmt.Sprintf("%#v", field)).
			Msg("getting field...")
		// Get the name of the field.
		fieldName := structType.Field(i).Name
		configEmptyLogger.Debug().
			Str("fieldName", fieldName).
			Msg("getting field name...")
		// Set the path of the field.
		fieldPath := path + "." + fieldName
		configEmptyLogger.Debug().
			Str("fieldPath", fieldPath).
			Msg("setting field path...")
		// Switch fild type.
		switch field.Kind() {
		// Check if the field is a struct.
		case reflect.Struct:
			if err := ConfigIsEmpty(field.Interface(), fieldPath, logger); err != nil {
				logger.Error().
					Err(err).
					Msg("")
				return err
			}
		// Check if the field is a pointer.
		case reflect.Ptr:
			// Check if the field is nil.
			if field.IsNil() {
				err := errors.New(fmt.Sprintf("%s is nil", fieldPath))
				configEmptyLogger.Error().
					Err(err).
					Msg("")
				return err
			} else {
				// Check if the field is empty.
				if field.Elem().Kind() == reflect.Struct {
					if err := ConfigIsEmpty(field.Elem().Interface(), fieldPath, logger); err != nil {
						configEmptyLogger.Error().
							Err(err).
							Msg("")
						return err
					}
				} else {
					// Check if the field is zero.
					if field.Elem().IsZero() {
						err := errors.New(fmt.Sprintf("%s is empty", fieldPath))
						configEmptyLogger.Error().
							Err(err).
							Msg("")
						return err
					}
					configEmptyLogger.Debug().
						Str("fieldPath", fieldPath).
						Msg(fmt.Sprintf("%s: not empty\n", fieldPath))
				}
			}
		// Check if the field is zero.
		default:
			if field.Interface() == reflect.Zero(field.Type()).Interface() {
				err := errors.New(fmt.Sprintf("%s is empty", fieldPath))
				configEmptyLogger.Error().
					Err(err).
					Msg("")
				return err
			}
			configEmptyLogger.Debug().
				Str("fieldPath", fieldPath).
				Msg(fmt.Sprintf("%s: not empty\n", fieldPath))
		}
	}

	configEmptyLogger.Debug().
		Str("path", path).
		Msg("checking if config is empty successful")

	return nil
}

func LoadConfig(path string, config interface{}) error {
	// Add new default fields to the logger for this function.
	loadConfigLogger := log.With().Str("func", "LoadConfig").Logger()

	loadConfigLogger.Debug().
		Str("path", path).
		Str("config", fmt.Sprintf("%#v", config)).
		Msg("loading config")

	// Set the config path, name and type.
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Read in the config file.
	err := viper.ReadInConfig()
	if err != nil {
		loadConfigLogger.Error().
			Err(err).
			Msg("failed to read config")
		return err
	}

	// Unmarshal the config file.
	err = viper.Unmarshal(&config)
	if err != nil {
		loadConfigLogger.Error().
			Err(err).
			Msg("failed to unmarshal config")
		return err
	}

	loadConfigLogger.Debug().
		Str("config", fmt.Sprintf("%#v", &config)).
		Msg("loading config success")

	return nil
}
