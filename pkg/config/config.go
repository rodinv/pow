package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

// Server is a server config
type Server struct {
	Port string `env:"POW_PORT" env-default:"8081"`
	Host string `env:"POW_HOST" env-default:"localhost"`
	Bits int64  `env:"POW_BITS" env-default:"24"`
}

// Get gets config from env
func Get() (*Server, error) {
	cfg := Server{}

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
