package config

import (
	"github.com/caarlos0/env/v6"
	"log"
)

func Load(config interface{}) error {
	err := env.Parse(config)
	log.Printf("Config loaded: %v", config)
	return err
}
