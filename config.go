package main

import (
	"encoding/json"
	"log"
	"os"
)

// Config struct holds all our configuration
type Config struct {
	WavesNode       string `json:"waves_node"`
	WavesNodeAPIKey string `json:"waves_node_api_key"`
	NodeAddress     string `json:"node_address"`
	Debug           bool   `json:"debug"`
	SSL             bool   `json:"ssl"`
	TelegramAPIKey  string `json:"telegram_api_key"`
	EmailAddress    string `json:"email_address"`
	TokenID         string `json:"token_id"`
	Hostname        string `json:"hostname"`
	PricesUrl       string `json:"prices_url"`
}

// Load method loads configuration file to Config struct
func (sc *Config) Load(configFile string) error {
	file, err := os.Open(configFile)

	if err != nil {
		log.Printf("[Config.Load] Got error while opening config file: %v", err)
		return err
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&sc)

	if err != nil {
		log.Printf("[Config.Load] Error while decoding JSON: %v", err)
		return err
	}

	return nil
}

func initConfig() *Config {
	c := &Config{}
	c.Load("config.json")
	return c
}
