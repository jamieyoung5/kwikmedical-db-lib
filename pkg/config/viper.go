package config

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"reflect"
	"time"
)

// EnvVarPrefix establishes that, by default, our environment variable naming prefix is "APP_".
// Note: Once we eliminate use of YAML files, we have to do that manually since Viper only
// applies the 'configName' automatically when using both env vars plus a config file.
const EnvVarPrefix = "APP_"
const DefaultEnvVarPrefix = "APP"

type AppViper struct {
	*viper.Viper
	prefix string
}

// NewAppViper constructs a Viper config using the default prefix that we will then extend with additional methods.
func NewAppViper() *AppViper {
	av := &AppViper{
		Viper:  newViper(),
		prefix: DefaultEnvVarPrefix,
	}
	return av
}

// NewAppViperWithPrefix constructs a Viper config using a custom prefix that we will then extend with additional methods.
func NewAppViperWithPrefix(prefix string) *AppViper {
	av := &AppViper{
		Viper:  newViper(),
		prefix: prefix,
	}
	return av
}

func newViper() *viper.Viper {
	viperConfig := viper.New()
	viperConfig.AllowEmptyEnv(true)
	_ = viperConfig.BindPFlags(pflag.CommandLine)
	return viperConfig
}

func (av *AppViper) SetDefaults(keysAndValues *map[string]interface{}) {
	for k, v := range *keysAndValues {
		av.SetDefault(k, v)
	}
}

func (av *AppViper) BindVariables(keys []string) {
	for _, k := range keys {
		_ = av.BindEnv(k)
	}
}

func (av *AppViper) SetAndBindDefaults(keysAndValues *map[string]interface{}) {
	for k, v := range *keysAndValues {
		av.setAndBindDefault(k, v)
	}
}

func (av *AppViper) setAndBindDefault(key string, value interface{}) {
	av.SetDefault(key, value)
	_ = av.BindEnv(key)
}

func GetConfig[T interface{}](defaultConfig T) T {
	return GetConfigUsingPrefix(defaultConfig, DefaultEnvVarPrefix)
}

func GetConfigUsingPrefix[T interface{}](defaultConfig T, prefix string) T {
	av := NewAppViperWithPrefix(prefix)

	typeOfConfig := reflect.TypeOf(defaultConfig)

	defaultStructValue := reflect.ValueOf(defaultConfig)

	for i := 0; i < defaultStructValue.NumField(); i++ {
		key := typeOfConfig.Field(i).Name
		key = av.getKeyWithPrefix(key)

		defaultValue := defaultStructValue.Field(i).Interface()
		av.setAndBindDefault(key, defaultValue)
	}

	result := defaultConfig

	resultValue := reflect.ValueOf(result)

	for i := 0; i < resultValue.NumField(); i++ {
		resultElem := reflect.ValueOf(&result).Elem()
		field := resultElem.Field(i)
		key := typeOfConfig.Field(i).Name
		key = av.getKeyWithPrefix(key)
		switch field.Interface().(type) {
		case string:
			envVarValue := av.GetString(key)
			field.SetString(envVarValue)
			break
		case []string:
			envVarValue := av.GetStringSlice(key)
			field.Set(reflect.ValueOf(envVarValue))
			break
		case int64:
			envVarValue := av.GetInt64(key)
			field.SetInt(envVarValue)
			break
		case int:
			envVarValue := av.GetInt(key)
			field.SetInt(int64(envVarValue))
			break
		case bool:
			envVarValue := av.GetBool(key)
			field.SetBool(envVarValue)
			break
		case time.Duration:
			envVarValue := av.GetDuration(key)
			field.Set(reflect.ValueOf(envVarValue))
		}
	}

	return result
}

// GetExpectedEnvVarNames is a helper method to list the expected environment variable names
// that are needed to successfully map onto a given config struct
func GetExpectedEnvVarNames[T interface{}](prefix string) []string {
	av := NewAppViperWithPrefix(prefix)

	var configInstance T
	typeOfConfig := reflect.TypeOf(configInstance)

	var result []string
	for i := 0; i < typeOfConfig.NumField(); i++ {
		key := typeOfConfig.Field(i).Name
		key = av.getKeyWithPrefix(key)
		result = append(result, key)
	}

	return result
}

func (av *AppViper) getKeyWithPrefix(key string) string {
	key = strcase.ToScreamingSnake(key)
	return fmt.Sprintf("%v_%v", av.prefix, key)
}
