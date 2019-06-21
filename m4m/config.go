package m4m

import (
	"encoding/json"
	"os"
)

type Config struct {
	Stripe StripeConfig
	FabMan FabManConfig
}

type StripeConfig struct {
	Secret        string
	SigningSecret string
}

type FabManConfig struct {
	Account int
	Token   string
}

func GetConfig() *Config {
	configFile, err := os.Open("config.json")
	if err != nil {
		panic(err)
	}
	cfg := &Config{}

	if err := json.NewDecoder(configFile).Decode(cfg); err != nil {
		panic(err)
	}

	return cfg
}
