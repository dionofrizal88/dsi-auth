package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Configuration struct is used when get configuration.
type Configuration struct {
	AppEnv     string `json:"APP_ENV"`
	AppName    string `json:"APP_NAME"`
	AppPort    string `json:"APP_PORT"`
	AppSecret  string `json:"APP_SECRET"`
	DBUsername string `json:"DB_USERNAME"`
	DBPassword string `json:"DB_PASSWORD"`
	DBHost     string `json:"DB_HOST"`
	DBPort     string `json:"DB_PORT"`
	DBName     string `json:"DB_NAME"`
	TestDBName string `json:"TEST_DB_NAME"`
	RedisHost  string `json:"REDIS_HOST"`
	RedisPort  string `json:"REDIS_PORT"`
	RedisDB    int    `json:"REDIS_DB"`
}

// GetConfig is a function to get configuration.
func GetConfig(path string) Configuration {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening config file: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	conf := Configuration{}
	err = decoder.Decode(&conf)
	if err != nil {
		fmt.Printf("Error decoding JSON config: %v", err)
	}

	return conf
}
