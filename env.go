// Package env provides a generic interface for loading environment variables
// from a config file and unmarshaling them into a struct.
package env

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/VinukaThejana/go-utils/logger"
	"github.com/spf13/viper"
)

// Env is an interface that defines the methods for loading environment variables.
type Env interface {
	Load(path ...string)
}

// Load loads environment variables from the given path and unmarshals them into the given struct.
func Load[T any](e *T, path ...string) {
	configPath := "."
	configFile := ".env"

	v := viper.New()

	if len(path) > 2 {
		logger.Errorf(fmt.Errorf("invalid set of parameters are provided"))
	}

	if len(path) > 0 {
		if len(path) == 2 {
			configFile = path[1]
		}
		configPath = path[0]

		if strings.HasSuffix(path[0], "/") {
			configFile = fmt.Sprintf("%s%s", configPath, configFile)
		} else {
			configFile = fmt.Sprintf("%s/%s", configPath, configFile)
		}
	}

	_, err := os.Stat(configFile)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			lf(err)
		}

		lf(parseEnvVars(environ(), e))
	} else {
		v.AddConfigPath(configPath)
		v.SetConfigFile(configFile)

		lf(v.ReadInConfig())
		lf(v.Unmarshal(e))
	}

	logger.Validatef(e)
}

// environ returns a map of environment variables and their values.
func environ() map[string]string {
	m := make(map[string]string)
	for _, s := range os.Environ() {
		a := strings.Split(s, "=")
		m[a[0]] = a[1]
	}

	return m
}

// parseEnvVars parses the environment variables in the given map and unmarshals them into the given struct.
func parseEnvVars[T any](envMap map[string]string, e *T) error {
	objValue := reflect.ValueOf(e).Elem()
	objType := objValue.Type()

	for i := 0; i < objType.NumField(); i++ {
		field := objType.Field(i)
		envKey := field.Tag.Get("mapstructure")
		envValue, ok := envMap[envKey]
		if !ok {
			continue
		}

		fieldValue := objValue.Field(i)
		if !fieldValue.CanSet() {
			return fmt.Errorf("field %s is not settable", field.Name)
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(envValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(envValue, 10, 64)
			if err != nil {
				return fmt.Errorf("failed to parse %s as int: %v", envKey, err)
			}

			fieldValue.SetInt(val)
		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(envValue, 64)
			if err != nil {
				return fmt.Errorf("failed to parse %s as float: %v", envKey, err)
			}

			fieldValue.SetFloat(val)
		case reflect.Bool:
			val, err := strconv.ParseBool(envValue)
			if err != nil {
				return fmt.Errorf("failed to parse %s as bool: %v", envKey, err)
			}

			fieldValue.SetBool(val)
		default:
			if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
				val, err := time.Parse(time.RFC3339, envValue)
				if err != nil {
					return fmt.Errorf("failed to parse %s as time.Time: %v", envKey, err)
				}

				fieldValue.Set(reflect.ValueOf(val))
			} else {
				return fmt.Errorf("unsupported type for field %s", field.Name)
			}
		}
	}

	return nil
}

// lf logs the error and exits the program.
func lf(err error) {
	if err != nil {
		logger.Errorf(err)
	}
}
