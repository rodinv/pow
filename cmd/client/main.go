package main

import (
	"fmt"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/rodinv/pow/pkg/client"
	"github.com/rodinv/pow/pkg/config"
)

func main() {
	// init logger
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// init config
	cfg, err := config.Get()
	if err != nil {
		log.Fatal().Err(err).Msg("initializing config")
		return
	}

	// get data and print
	cli := client.New(cfg.Host, cfg.Port)
	quote, err := cli.GetQuote()
	if err != nil {
		log.Fatal().Err(err).Msg("getting quote")
		return
	}

	fmt.Println(quote)
}
