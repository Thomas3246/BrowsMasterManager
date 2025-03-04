package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	BotToken    string `json:"botToken"`
	MasterPhone string `json:"masterPhone"`
	AboutMaster string `json:"aboutMaster"`
}

func NewConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, err
}
