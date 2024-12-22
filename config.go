package main

import (
	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	DatabaseURL      string `env:"DATABASE_URL,required"`  // MySQL DSN
	DiscordToken     string `env:"DISCORD_TOKEN,required"` // Discord Bot Token
	ServerPort       int    `env:"SERVER_PORT" envDefault:"8080"`
	RotatorDelay     int    `env:"ROTATOR_DELAY" envDefault:"1"`     // in hours
	RotatorThreshold int    `env:"ROTATOR_THRESHOLD" envDefault:"1"` // in hours
}

func LoadConfig() (Config, error) {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return cfg, err
	}

	log.Info().Msgf("Loaded config: %+v", cfg)
	return cfg, nil
}
