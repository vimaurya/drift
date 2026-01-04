package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	DatabaseURL string `json:"database_url"`
	Dir string `json:"migration_dir"`
}

func Save(cfg Config) error {
	data, err := json.MarshalIndent(cfg, "", " ")
	if err!=nil{
		return err
	}
	return os.WriteFile(".gomigrate.json", data, 0644)
}

func Load() (*Config, error) {
	data, err := os.ReadFile(".gomigrate.json")
	if err!=nil{
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	
	return &config, err
}
