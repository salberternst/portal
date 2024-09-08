package utils

import (
	"os"
	"strconv"
	"sync"
)

type Config struct {
	EnableDeviceApi     bool
	EnableFusekiBackend bool
	EnableUsers         bool
}

var (
	config *Config
	once   sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		config = &Config{
			EnableDeviceApi:     ReadBoolEnvVariable("ENABLE_DEVICE_API"),
			EnableFusekiBackend: ReadBoolEnvVariable("ENABLE_FUSEKI_BACKEND"),
			EnableUsers:         ReadBoolEnvVariable("ENABLE_USERS"),
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
