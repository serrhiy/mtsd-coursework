package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
)

type RequestFormat struct {
	service, method string
	data            any
}

func Call(function any, args ...any) error {
	value := reflect.ValueOf(function)
	if value.Kind() != reflect.Func {
		return errors.New("value is not a function")
	}
	if value.Type().NumIn() != len(args) {
		return errors.New("invalid arguments count")
	}
	arguments := make([]reflect.Value, value.Type().NumIn())
	for index := range args {
		functionArgumentType := value.Type().In(index)
		argumentValue := reflect.ValueOf(args[index])
		if !argumentValue.Type().ConvertibleTo(functionArgumentType) {
			return errors.New("invalid arguments format")
		}
		converted := argumentValue.Convert(functionArgumentType)
		arguments[index] = converted
	}
	value.Call(arguments)
	return nil
}

func ReadFields(json map[string]any, fields []string) (map[string]any, error) {
	response := make(map[string]any)
	for _, field := range fields {
		value, ok := json[field]
		if !ok {
			return nil, fmt.Errorf("key '%s' is absent", field)
		}
		response[field] = value
	}
	return response, nil
}

func ParseRequestJson(json map[string]any) (RequestFormat, error) {
	packet, err := ReadFields(json, []string{"service", "method"})
	if err != nil {
		return RequestFormat{}, err
	}
	service, ok := packet["service"].(string)
	if !ok {
		return RequestFormat{}, errors.New("invalid service key")
	}
	method, ok := packet["method"].(string)
	if !ok {
		return RequestFormat{}, errors.New("invalid method key")
	}
	return RequestFormat{service: service, method: method, data: json["data"]}, nil
}

func GetConfig(filepath string) (Config, error) {
	content, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(content, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func DatabaseURL(config DatabaseConfig) string {
	schema := config.Engine + "://"
	credentials := config.User + ":" + config.Password
	destination := config.Host + ":" + config.Port
	return schema + credentials + "@" + destination + "/" + config.Database
}
