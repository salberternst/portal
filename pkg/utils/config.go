package utils

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	EnableThingsboard bool
	EnableFuseki      bool
	EnableUsers       bool
}

var (
	config *Config
	once   sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		config = &Config{
			EnableThingsboard: ReadBoolEnvVariable("ENABLE_THINGSBOARD"),
			EnableFuseki:      ReadBoolEnvVariable("ENABLE_FUSEKI"),
			EnableUsers:       ReadBoolEnvVariable("ENABLE_USERS"),
		}
	})
	return config
}

func ReadBoolEnvVariable(key string) bool {
	value := os.Getenv(key)
	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return false
	}
	return boolValue
}
