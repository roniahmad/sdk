package sdk

import (
	"errors"
	"os"
	"reflect"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/xhit/go-str2duration/v2"
)

type AppConfig[T any] interface {
	GetConfig() T
}

type configHolder[T any] struct {
	config T
}

// Create a new configuration object and return `AppConfig[T]` interface.
//
// Function `beforeConfigLoadedHook` and `afterConfigLoadedHook` can be used to modify configuration object.
//
// For example:
//
//	 cfg := configloader.New("./config.yaml",
//		 // before config loaded hook, can be used to load default value or load multiple config files
//		 func(c appconfig.Config, loader *viper.Viper) {
//			 c.LoadDefaultValue(loader)
//			 c.LoadAccessControlListConfig(loader)
//		 },
//		 // after config loaded hooks, can be used  to validate the configuration
//		 func(c appconfig.Config, _ *viper.Viper) {
//			 if err := c.Validate(); err != nil {
//				 log.Fatalf("Error validating configuration: %s", err)
//			 }
//		 },
//	 )
//	 config := cfg.GetConfig()  // get configuration

func New[T any](
	configPath string,
	beforeConfigLoadedHook func(config *T, loader *viper.Viper),
	afterConfigLoadedHook func(config *T, loader *viper.Viper),
) AppConfig[T] {
	configStruct := new(T)

	if _, err := os.Stat(configPath); errors.Is(err, os.ErrNotExist) {
		log.Fatal().Msg("Configuration file not exist")
	}

	appConfigLoader := viper.New()
	appConfigLoader.SetConfigFile(configPath)
	appConfigLoader.AutomaticEnv()

	if err := appConfigLoader.ReadInConfig(); err != nil {
		log.Fatal().Msgf("Failed reading configuration: %s", err)
	}

	// before config loaded hook
	if beforeConfigLoadedHook != nil {
		beforeConfigLoadedHook(configStruct, appConfigLoader)
	}

	// replace env placeholder with environment variables
	for _, key := range appConfigLoader.AllKeys() {
		value := appConfigLoader.Get(key)

		switch parsedValue := value.(type) {
		case string:
			appConfigLoader.Set(key, os.ExpandEnv(parsedValue))
		default:
			appConfigLoader.Set(key, parsedValue)
		}
	}

	// marshall configuration yaml
	if err := appConfigLoader.Unmarshal(&configStruct, func(m *mapstructure.DecoderConfig) {
		m.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			stringToTimeDurationHookFunc(),
			mapstructure.RecursiveStructToMapHookFunc(),
			mapstructure.TextUnmarshallerHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)
	}); err != nil {
		log.Fatal().Msgf("Error unmarshaling configuration: %s", err)
	}

	// after config loaded hooks
	if afterConfigLoadedHook != nil {
		afterConfigLoadedHook(configStruct, appConfigLoader)
	}

	return &configHolder[T]{
		config: *configStruct,
	}
}

// Get the configuration struct
func (c *configHolder[T]) GetConfig() T {
	return c.config
}

// stringToTimeDurationHookFunc returns a DecodeHookFunc that converts
// strings to time.Duration.
func stringToTimeDurationHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.String {
			return data, nil
		}
		if t != reflect.TypeOf(time.Duration(5)) {
			return data, nil
		}

		// Convert it by parsing
		return str2duration.ParseDuration(data.(string))
	}
}
