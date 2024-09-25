package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func load(cfg interface{}) error {
	v := reflect.ValueOf(cfg).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			continue
		}

		parts := strings.Split(envTag, ",")
		envVar := parts[0]
		var defaultValue string
		var isRequired bool

		for _, part := range parts[1:] {
			if part == "required" {
				isRequired = true
			} else if strings.HasPrefix(part, "default=") {
				defaultValue = part[8:]
			}
		}

		value, exists := os.LookupEnv(envVar)
		if !exists {
			if defaultValue != "" {
				value = defaultValue
			} else if isRequired {
				return fmt.Errorf("environment variable %s is required but not set", envVar)
			}
		}

		// Set the field value
		fieldValue := v.Field(i)
		if !fieldValue.CanSet() {
			return fmt.Errorf("cannot set field %s", field.Name)
		}

		switch fieldValue.Kind() {
		case reflect.String:
			fieldValue.SetString(value)
		case reflect.Int:
			intValue, err := strconv.Atoi(value)
			if err != nil {
				return fmt.Errorf("error converting %s to int: %v", value, err)
			}
			fieldValue.SetInt(int64(intValue))
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("error converting %s to bool: %v", value, err)
			}
			fieldValue.SetBool(boolValue)
		default:
			return fmt.Errorf("unsupported field type %s", field.Type)
		}
	}

	return nil
}
