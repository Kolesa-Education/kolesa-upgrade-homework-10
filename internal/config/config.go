package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Bot struct {
		Token string `json:"token"`
	} `json:"bot"`
	Gorm struct {
		Dsn string `json:"dsn"`
	} `json:"gorm"`
}

func GetConfig(path string) (*Config, error) {
	var cfg Config
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
