package config

import "time"

type Config struct {
	JsonAPIConfig ProviderConfig
	MockAPIConfig ProviderConfig
	PollInterval  time.Duration
}

type ProviderConfig struct {
	Endpoint string
	APIKey   string
}

func LoadConfig() (*Config, error) {
	return &Config{
		JsonAPIConfig: ProviderConfig{
			Endpoint: "http://localhost:3000",
			APIKey:   "",
		},
		MockAPIConfig: ProviderConfig{
			Endpoint: "https://66e63e3917055714e58930b7.mockapi.io/mockapi",
			APIKey:   "",
		},
		PollInterval: 10 * time.Second,
	}, nil
}
